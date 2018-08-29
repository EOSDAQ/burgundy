// Package service ...
//
// This layer will act as the business process handler.
// Any process will handled here. This layer will decide, which repository layer will use.
// And have responsibility to provide data to serve into delivery.
// Process the data doing calculation or anything will done here.
//
// Service layer will accept any input from Delivery layer,
// that already sanitized, then process the input could be storing into DB ,
// or Fetching from DB ,etc.
//
// This Service layer will depends to Repository Layer
package service

import (
	"burgundy/models"
	"burgundy/util"
	"context"

	"go.uber.org/zap"
)

var mlog *zap.SugaredLogger

func init() {
	mlog, _ = util.InitLog("service", "console")
}

// UserService ...
type UserService interface {
	GetByID(ctx context.Context, accountName string) (*models.User, error)
	Store(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, accountName string) (bool, error)

	// Login
	Login(ctx context.Context, accountName, accountHash string) (*models.User, error)

	// Email
	ConfirmEmail(ctx context.Context, accountName, email, emailHash string) (*models.User, error)
	RevokeEmail(ctx context.Context, accountName, email, emailHash string) (*models.User, error)

	// OTP
	GenerateOTPKey(ctx context.Context, accountName string) (string, error)
	RevokeOTP(ctx context.Context, accountName string) error
	ValidateOTP(ctx context.Context, accountName, code string) (bool, error)
}

// EosdaqService ...
type EosdaqService interface {
	UpdateOrderbook(ctx context.Context, obs []*models.OrderBook, orderType models.OrderType) error
	UpdateTransaction(ctx context.Context, txs []*models.EosdaqTx) error
}
