package service

import (
	"burgundy/models"
	"burgundy/repository"
	"context"
	"time"
)

// UserService ...
type UserService interface {
	GetByID(ctx context.Context, accountName string) (*models.User, error)
	Store(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, accountName string) (bool, error)
}

type userUsecase struct {
	userRepo   repository.UserRepository
	ctxTimeout time.Duration
}

// NewUserService ...
func NewUserService(ur repository.UserRepository,
	timeout time.Duration) UserService {
	return &userUsecase{
		userRepo:   ur,
		ctxTimeout: timeout,
	}
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

	return uuc.userRepo.Store(innerCtx, user)
}

// Delete ...
func (uuc userUsecase) Delete(ctx context.Context, accountName string) (result bool, err error) {
	innerCtx, cancel := context.WithTimeout(ctx, uuc.ctxTimeout)
	defer cancel()

	return uuc.userRepo.Delete(innerCtx, accountName)
}
