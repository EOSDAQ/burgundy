package service

import (
	"burgundy/conf"
	"burgundy/models"
	"burgundy/repository"
	"context"
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

	if len(obs) == 0 {
		return nil
	}

	// get db old
	orderBooks, err := eu.eosdaqRepo.GetOrderBook(innerCtx, obs[0].Type)
	if err != nil {
		mlog.Errorw("UpdateOrderbook get", "contract", eu.contract, "err", err)
		return err
	}
	mlog.Debugw("UpdateOrderbook db read", "cont", eu.contract, "data", orderBooks)
	orderMaps := make(map[uint]*models.OrderBook)
	for _, o := range orderBooks {
		orderMaps[o.ID] = o
	}
	// diff obs,db
	addBooks := []*models.OrderBook{}
	for _, n := range obs {
		if _, ok := orderMaps[n.ID]; !ok {
			addBooks = append(addBooks, n)
		} else if ok {
			delete(orderMaps, n.ID)
		}
	}
	mlog.Debugw("UpdateOrderbook db add", "cont", eu.contract, "data", addBooks)

	// insert collection
	if err = eu.eosdaqRepo.SaveOrderBook(innerCtx, addBooks); err != nil {
		mlog.Errorw("UpdateOrderbook save", "contract", eu.contract, "err", err, "add", addBooks)
		return err
	}
	delBooks := []*models.OrderBook{}
	for _, d := range orderMaps {
		delBooks = append(delBooks, d)
	}
	mlog.Debugw("UpdateOrderbook db del", "cont", eu.contract, "data", delBooks)
	// delete collection
	if err = eu.eosdaqRepo.DeleteOrderBook(innerCtx, delBooks); err != nil {
		mlog.Errorw("UpdateOrderbook delete", "contract", eu.contract, "err", err, "del", delBooks)
		return err
	}

	// websocket broadcast

	return
}

// UpdateTransaction ...
func (eu eosdaqUsecase) UpdateTransaction(ctx context.Context, txs []*models.EosdaqTx) (err error) {
	innerCtx, cancel := context.WithTimeout(ctx, eu.ctxTimeout)
	defer cancel()

	if len(txs) == 0 {
		return nil
	}

	// get db old
	dbtxs, err := eu.eosdaqRepo.GetTransactions(innerCtx, txs)
	if err != nil && err.Error() != "record not found" {
		mlog.Errorw("UpdateTransactions get", "contract", eu.contract, "err", err)
		return err
	}
	mlog.Debugw("UpdateTransaction db read", "cont", eu.contract, "data", dbtxs)
	txMaps := make(map[uint]struct{})
	for _, t := range dbtxs {
		txMaps[t.ID] = struct{}{}
	}

	// diff txs,db
	addtxs := []*models.EosdaqTx{}
	for _, t := range txs {
		if _, ok := txMaps[t.ID]; !ok {
			addtxs = append(addtxs, t)
		}
	}
	mlog.Debugw("UpdateTransactions db add", "cont", eu.contract, "data", addtxs)

	if err = eu.eosdaqRepo.SaveTransaction(innerCtx, addtxs); err != nil {
		mlog.Errorw("UpdateTransaction", "contract", eu.contract, "txs", addtxs, "err", err)
		return err
	}

	return nil
}
