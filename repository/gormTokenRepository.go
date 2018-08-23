package repository

import (
	"burgundy/conf"
	models "burgundy/models"
	"context"

	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
)

type gormTokenRepository struct {
	Conn *gorm.DB
}

// NewGormTokenRepository ...
func NewGormTokenRepository(burgundy conf.ViperConfig, Conn *gorm.DB) TokenRepository {
	Conn = Conn.AutoMigrate(&models.Token{})
	g := &gormTokenRepository{Conn}
	for _, t := range models.TokenInit(burgundy.GetString("eos_baseSymbol")) {
		g.createToken(context.Background(), t)
	}
	return g
}

func (g *gormTokenRepository) GetTokens(ctx context.Context) (ts []*models.Token, err error) {
	scope := g.Conn.Find(&ts)
	if scope.RowsAffected == 0 {
		return nil, nil
	}
	return ts, scope.Error
}

func (g *gormTokenRepository) GetToken(ctx context.Context, symbol string) (token *models.Token, err error) {
	scope := g.Conn.New()
	scope.Where(models.Token{Symbol: symbol}).First(&token)
	if scope.RowsAffected == 0 {
		return nil, nil
	}
	return token, scope.Error
}

func (g *gormTokenRepository) createToken(ctx context.Context, token *models.Token) (err error) {
	g.Conn.Where(models.Token{Symbol: token.Symbol}).FirstOrCreate(token)
	if g.Conn.Error != nil {
		mlog.Errorw("UpdateToken", "err", g.Conn.Error)
		return errors.Annotatef(g.Conn.Error, "UpdateToken error [%s]", token.Symbol)
	}
	return nil
}

func (g *gormTokenRepository) UpdateToken(ctx context.Context, token *models.Token) (err error) {
	g.Conn.Where(models.Token{Symbol: token.Symbol}).Updates(token)
	if g.Conn.Error != nil {
		mlog.Errorw("UpdateToken", "err", g.Conn.Error)
		return errors.Annotatef(g.Conn.Error, "UpdateToken error [%s]", token.Symbol)
	}
	return nil
}
