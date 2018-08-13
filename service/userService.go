package service

import (
	"burgundy/conf"
	"burgundy/eosdaq"
	"burgundy/models"
	"burgundy/repository"
	"context"
	"time"

	"github.com/juju/errors"
)

type userUsecase struct {
	userRepo   repository.UserRepository
	eosAPI     *eosdaq.EosdaqAPI
	ctxTimeout time.Duration
}

// NewUserService ...
func NewUserService(burgundy conf.ViperConfig,
	ur repository.UserRepository,
	timeout time.Duration) (UserService, error) {
	eosapi, err := eosdaq.NewAPI(burgundy, eosdaq.NewEosnet(
		burgundy.GetString("eos_host"),
		burgundy.GetInt("eos_port"),
		burgundy.GetString("eos_acctcontract"),
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

	return uuc.userRepo.GetByID(innerCtx, accountName)
}

// Store ...
func (uuc userUsecase) Store(ctx context.Context, user *models.User) (u *models.User, err error) {
	innerCtx, cancel := context.WithTimeout(ctx, uuc.ctxTimeout)
	defer cancel()

	if err = uuc.eosAPI.DoAction(uuc.eosAPI.RegisterAction(user.AccountName)); err != nil {
		mlog.Infow("userUsecase register error", "user", user, "err", err)
		// To pass Already Exist
	}

	u, err = uuc.userRepo.Store(innerCtx, user)
	if err != nil {
		// To pass Already Exist
		/*
			mlog.Infow("userUsecase Store error", "user", user, "err", err)
			rbErr := uuc.eosAPI.UnregisterUser(user.AccountName)
			if rbErr != nil {
				mlog.Infow("userUsecase register rollback error", "user", user, "err", rbErr)
			}
			return nil, errors.Annotatef(err, "user[%v]", user)
		*/
	}

	return
}

// Update ...
func (uuc userUsecase) Update(ctx context.Context, user *models.User) (u *models.User, err error) {
	innerCtx, cancel := context.WithTimeout(ctx, uuc.ctxTimeout)
	defer cancel()

	return uuc.userRepo.Update(innerCtx, user)
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
