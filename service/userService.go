package service

import (
	"burgundy/conf"
	"burgundy/eosdaq"
	"burgundy/models"
	"burgundy/repository"
	"context"
	"time"

	eos "github.com/eoscanada/eos-go"
	"github.com/juju/errors"
)

type userUsecase struct {
	userRepo   repository.UserRepository
	eosAPI     *eosdaq.API
	ctxTimeout time.Duration
}

// NewUserService ...
func NewUserService(burgundy *conf.ViperConfig,
	ur repository.UserRepository,
	timeout time.Duration) (UserService, error) {
	eosapi, err := eosdaq.NewAPI(burgundy, eosdaq.NewEosnet(
		burgundy.GetString("eos_host"),
		burgundy.GetInt("eos_port"),
		burgundy.GetString("eos_acctcontract"),
		burgundy.GetString("eos_managecontract"),
	))
	if err != nil {
		return nil, errors.Annotatef(err, "NewUserService")
	}
	return &userUsecase{
		userRepo:   ur,
		eosAPI:     eosapi,
		ctxTimeout: timeout,
	}, nil
}

// GetByID ...
func (uuc userUsecase) GetByID(ctx context.Context, accountName string) (u *models.User, err error) {
	innerCtx, cancel := context.WithTimeout(ctx, uuc.ctxTimeout)
	defer cancel()

	u, err = uuc.userRepo.GetByID(innerCtx, accountName)
	if err != nil {
		return nil, err
	}

	u.EmailHash = nil

	return
}

// Store ...
func (uuc userUsecase) Store(ctx context.Context, user *models.User) (u *models.User, err error) {
	innerCtx, cancel := context.WithTimeout(ctx, uuc.ctxTimeout)
	defer cancel()

	u, err = uuc.userRepo.Store(innerCtx, user)
	if err != nil {
		return nil, err
	}

	u.EmailHash = nil

	return
}

// Delete ...
func (uuc userUsecase) Delete(ctx context.Context, accountName string) (result bool, err error) {
	innerCtx, cancel := context.WithTimeout(ctx, uuc.ctxTimeout)
	defer cancel()

	var eoserr, dberr error
	if err = uuc.eosAPI.DoAction(uuc.eosAPI.UnregisterAction(accountName)); err != nil {
		mlog.Infow("userUsecase unregister error", "user", accountName, "err", err)
		eoserr = errors.Annotatef(err, "user[%s]", accountName)
	}

	result, err = uuc.userRepo.Delete(innerCtx, accountName)
	if err != nil {
		dberr = errors.Annotatef(err, "user[%s]", accountName)
	}

	if eoserr != nil && dberr != nil {
		result = false
	}

	return result, dberr
}

// Login ...
func (uuc userUsecase) Login(ctx context.Context, accName string) (ok bool, err error) {
	innerCtx, cancel := context.WithTimeout(ctx, uuc.ctxTimeout)
	defer cancel()

	u, err := uuc.userRepo.GetByID(innerCtx, accName)
	if err != nil {
		return false, errors.Annotatef(err, "Login GetByID[%s]", accName)
	}

	return u.AccountName == accName, nil
}

// ConfirmEmail ...
func (uuc userUsecase) ConfirmEmail(ctx context.Context, accName, email, emailHash string) (u *models.User, err error) {
	innerCtx, cancel := context.WithTimeout(ctx, uuc.ctxTimeout)
	defer cancel()

	dbuser, err := uuc.userRepo.GetByID(innerCtx, accName)
	if err != nil {
		return nil, errors.Annotatef(err, "ConfirmEmail GetByID[%s]", accName)
	}

	if !dbuser.ConfirmEmail(email, emailHash) {
		return nil, errors.NotValidf("ConfirmEmail Invalid Email Hash[%s]", emailHash)
	}
	uuc.updateContract(dbuser)

	u, err = uuc.userRepo.Update(innerCtx, dbuser)

	u.EmailHash = nil

	return
}

// RevokeEmail ...
func (uuc userUsecase) RevokeEmail(ctx context.Context, accName, email, emailHash string) (u *models.User, err error) {
	innerCtx, cancel := context.WithTimeout(ctx, uuc.ctxTimeout)
	defer cancel()

	dbuser, err := uuc.userRepo.GetByID(innerCtx, accName)
	if err != nil {
		return nil, errors.Annotatef(err, "RevokeEmail GetByID[%s]", accName)
	}

	dbuser.RevokeEmail(email, emailHash)
	uuc.updateContract(dbuser)

	u, err = uuc.userRepo.Update(innerCtx, dbuser)

	u.EmailHash = nil

	return
}

func (uuc userUsecase) updateContract(user *models.User) bool {
	var action *eos.Action
	if user.NeedRegister() {
		action = uuc.eosAPI.RegisterAction(user.AccountName)
	} else if user.NeedUnregister() {
		action = uuc.eosAPI.UnregisterAction(user.AccountName)
	}

	if action == nil {
		return false
	}

	if err := uuc.eosAPI.DoAction(action); err != nil {
		mlog.Infow("userUsecase register error", "user", user, "err", err)
		// To pass Already Exist
	}
	user.UpdateRegister()
	return true
}

func (uuc userUsecase) GenerateOTPKey(ctx context.Context, accountName string) (key string, err error) {
	innerCtx, cancel := context.WithTimeout(ctx, uuc.ctxTimeout)
	defer cancel()

	dbuser, err := uuc.userRepo.GetByID(innerCtx, accountName)
	if err != nil {
		return "", errors.Annotatef(err, "GetOTPKey GetByID error[%s]", accountName)
	}

	key, err = dbuser.GenerateOTPKey()
	if err != nil {
		return "", errors.Annotatef(err, "GetOTPKey GenerateOTPKey error[%s]", accountName)
	}
	_, err = uuc.userRepo.Update(innerCtx, dbuser)
	if err != nil {
		return "", errors.Annotatef(err, "GetOTPKey Update error[%s]", accountName)
	}

	return key, nil
}

func (uuc userUsecase) RevokeOTP(ctx context.Context, accountName string) (err error) {
	innerCtx, cancel := context.WithTimeout(ctx, uuc.ctxTimeout)
	defer cancel()

	dbuser, err := uuc.userRepo.GetByID(innerCtx, accountName)
	if err != nil {
		return errors.Annotatef(err, "RevokeOTP GetByID error[%s]", accountName)
	}

	dbuser.RemoveOTPKey()
	uuc.updateContract(dbuser)

	_, err = uuc.userRepo.Update(innerCtx, dbuser)
	if err != nil {
		return errors.Annotatef(err, "RevokeOTP Update error[%s]", accountName)
	}

	return nil
}

func (uuc userUsecase) ValidateOTP(ctx context.Context, accountName, code string) (ok bool, err error) {
	innerCtx, cancel := context.WithTimeout(ctx, uuc.ctxTimeout)
	defer cancel()

	dbuser, err := uuc.userRepo.GetByID(innerCtx, accountName)
	if err != nil {
		return false, errors.Annotatef(err, "ValidateOTP GetByID error[%s]", accountName)
	}

	if !dbuser.ValidateOTP(code) {
		return false, errors.NotValidf("ValidateOTP invalid code[%s]", code)
	}

	if uuc.updateContract(dbuser) {
		_, err = uuc.userRepo.Update(innerCtx, dbuser)
		if err != nil {
			return false, errors.Annotatef(err, "ValidateOTP Update error[%s]", accountName)
		}
	}

	return true, nil
}
