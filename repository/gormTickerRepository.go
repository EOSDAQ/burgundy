package repository

import (
	"burgundy/conf"
	models "burgundy/models"
	"context"

	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
)

type gormTickerRepository struct {
	Conn *gorm.DB
}

// NewGormTickerRepository ...
func NewGormTickerRepository(burgundy conf.ViperConfig, Conn *gorm.DB) TickerRepository {
	Conn = Conn.AutoMigrate(&models.Ticker{})
	g := &gormTickerRepository{Conn}
	for _, t := range models.TickerInit(burgundy.GetString("eos_baseSymbol")) {
		g.UpdateTicker(context.Background(), t)
	}
	return g
}

func (g *gormTickerRepository) GetTickers(ctx context.Context) (ts []*models.Ticker, err error) {
	scope := g.Conn.Find(&ts)
	if scope.RowsAffected == 0 {
		return nil, nil
	}
	return ts, scope.Error
}

func (g *gormTickerRepository) GetTicker(ctx context.Context, symbol string) (ticker *models.Ticker, err error) {
	scope := g.Conn.New()
	scope.Where(models.Ticker{TokenSymbol: symbol}).First(&ticker)
	if scope.RowsAffected == 0 {
		return nil, nil
	}
	return ticker, scope.Error
}

func (g *gormTickerRepository) UpdateTicker(ctx context.Context, ticker *models.Ticker) (err error) {

	mlog.Infow("UpdateTicker", "ticker", ticker)
	g.Conn.Debug().Where(models.Ticker{TokenSymbol: ticker.TokenSymbol}).FirstOrCreate(ticker)
	if g.Conn.Error != nil {
		mlog.Errorw("UpdateTicker", "err", g.Conn.Error)
		return errors.Annotatef(g.Conn.Error, "UpdateTicker error [%s]", ticker.TokenSymbol)
	}
	/*
		if scope.Error != nil {
			return errors.Annotatef(scope.Error, "UpdateTicker error [%s]", ticker.TokenSymbol)
		}
	*/
	return nil
}
