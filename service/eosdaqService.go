package service

import (
	"burgundy/conf"
	"burgundy/models"
	"burgundy/repository"
	"context"
	"time"
)

type eosdaqUsecase struct {
	eosdaqRepo repository.EosdaqRepository
	ctxTimeout time.Duration
}

// NewEosdaqService ...
func NewEosdaqService(burgundy conf.ViperConfig,
	er repository.EosdaqRepository,
	timeout time.Duration) (EosdaqService, error) {
	return &eosdaqUsecase{
		eosdaqRepo: er,
		ctxTimeout: timeout,
	}, nil
}

// UpdateOrderbook ...
func (eu eosdaqUsecase) UpdateOrderbook(ctx context.Context, obs []*models.OrderBook) (err error) {
	innerCtx, cancel := context.WithTimeout(ctx, eu.ctxTimeout)
	defer cancel()

	// obs new
	// get db old
	// diff obs,db
	// insert collection
	// delete collection
	// websocket broadcast

	return
}

// UpdateTransaction ...
func (eu eosdaqUsecase) UpdateTransaction(ctx context.Context, txs []*models.EosdaqTx) (err error) {
	innerCtx, cancel := context.WithTimeout(ctx, eu.ctxTimeout)
	defer cancel()

	return eu.eosdaqRepo.SaveTransaction(innerCtx, txs)
}
