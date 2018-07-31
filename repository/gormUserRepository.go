package repository

import (
	models "burgundy/models"
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
)

type gormUserRepository struct {
	Conn *gorm.DB
}

// NewGormUserRepository ...
func NewGormUserRepository(Conn *gorm.DB) UserRepository {
	Conn = Conn.AutoMigrate(&models.User{})
	return &gormUserRepository{Conn}
}

func (g *gormUserRepository) GetByID(ctx context.Context, accountName string) (user *models.User, err error) {
	user = &models.User{}
	scope := g.Conn.Where("account_name = ?", accountName).First(&user)
	if scope.Error != nil {
		return nil, scope.Error
	}

	if scope.RowsAffected == 0 {
		return nil, fmt.Errorf("record not found")
	}
	return user, nil

}

func (g *gormUserRepository) Update(ctx context.Context, user *models.User) (u *models.User, err error) {
	session := g.Conn.Where("account_name = ?", user.AccountName)
	if err = session.Update(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (g *gormUserRepository) Store(ctx context.Context, user *models.User) (*models.User, error) {
	scope := g.Conn.Where("account_name = ?", user.AccountName).FirstOrCreate(user)
	if err := scope.Error; err != nil {
		return nil, err
	}

	if scope.RowsAffected == 0 {
		if err := g.Conn.Create(user).Error; err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (g *gormUserRepository) Delete(ctx context.Context, accountName string) (bool, error) {
	if err := g.Conn.Where("account_name = ?", accountName).Delete(models.User{}).Error; err != nil {
		return false, err
	}
	return true, nil
}
