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
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, accountName string) (bool, error)
}
