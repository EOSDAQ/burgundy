package service

import (
	"burgundy/conf"
	"burgundy/models"
	"burgundy/repository"
	"context"
	"fmt"
	"time"
)

type eosdaqUsecase struct {
	contract   string
	eosdaqRepo repository.EosdaqRepository
	ctxTimeout time.Duration
}

// NewEosdaqService ...
func NewEosdaqService(burgundy conf.ViperConfig,
	contract string,
	er repository.EosdaqRepository,
	timeout time.Duration) (EosdaqService, error) {
	return &eosdaqUsecase{
		contract:   contract,
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
	for _, o := range obs {
		fmt.Printf("contract[%s] %v\n", eu.contract, o)
	}

	return
}

// UpdateTransaction ...
func (eu eosdaqUsecase) UpdateTransaction(ctx context.Context, txs []*models.EosdaqTx) (err error) {
	innerCtx, cancel := context.WithTimeout(ctx, eu.ctxTimeout)
	defer cancel()

	if err = eu.eosdaqRepo.SaveTransaction(innerCtx, txs); err != nil {
		mlog.Infow("UpdateTransaction", "contract", eu.contract, "txs", txs)
		return err
	}

	return nil
}
