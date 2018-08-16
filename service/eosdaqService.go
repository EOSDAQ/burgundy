package service

import (
	"burgundy/conf"
	"burgundy/models"
	"burgundy/repository"
	"context"
	"time"
)

type eosdaqUsecase struct {
	ticker     *models.Ticker
	eosdaqRepo repository.EosdaqRepository
	ctxTimeout time.Duration
}

// NewEosdaqService ...
func NewEosdaqService(burgundy conf.ViperConfig,
	t *models.Ticker,
	er repository.EosdaqRepository,
	timeout time.Duration) (EosdaqService, error) {
	return &eosdaqUsecase{
		ticker:     t,
		eosdaqRepo: er,
		ctxTimeout: timeout,
	}, nil
}

func (eu eosdaqUsecase) UpdateTicker(ctx context.Context, ticker *models.Ticker) (dbtick *models.Ticker, err error) {
	innerCtx, cancel := context.WithTimeout(ctx, eu.ctxTimeout)
	defer cancel()

	return eu.UpdateTicker(innerCtx, ticker)
}

// UpdateOrderbook ...
func (eu eosdaqUsecase) UpdateOrderbook(ctx context.Context, obs []*models.OrderBook, orderType models.OrderType) (err error) {
	innerCtx, cancel := context.WithTimeout(ctx, eu.ctxTimeout)
	defer cancel()

	// get db old
	orderBooks, err := eu.eosdaqRepo.GetOrderBook(innerCtx, orderType)
	if err != nil {
		mlog.Errorw("UpdateOrderbook get", "contract", eu.ticker.ContractAccount, "err", err)
		return err
	}
	//mlog.Debugw("UpdateOrderbook db read", "cont", eu.ticker.ContractAccount, "data", orderBooks)
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
	//mlog.Debugw("UpdateOrderbook db add", "cont", eu.ticker.ContractAccount, "data", addBooks)

	// insert collection
	if err = eu.eosdaqRepo.SaveOrderBook(innerCtx, addBooks); err != nil {
		mlog.Errorw("UpdateOrderbook save", "contract", eu.ticker.ContractAccount, "err", err, "add", addBooks)
		return err
	}
	delBooks := []*models.OrderBook{}
	for _, d := range orderMaps {
		delBooks = append(delBooks, d)
	}
	//mlog.Debugw("UpdateOrderbook db del", "cont", eu.ticker.ContractAccount, "data", delBooks)
	// delete collection
	if err = eu.eosdaqRepo.DeleteOrderBook(innerCtx, delBooks); err != nil {
		mlog.Errorw("UpdateOrderbook delete", "contract", eu.ticker.ContractAccount, "err", err, "del", delBooks)
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
		mlog.Errorw("UpdateTransactions get", "contract", eu.ticker.ContractAccount, "err", err)
		return err
	}
	//mlog.Debugw("UpdateTransaction db read", "cont", eu.ticker.ContractAccount, "data", dbtxs)
	txMaps := make(map[uint]struct{})
	for _, t := range dbtxs {
		txMaps[t.ID] = struct{}{}
	}

	// diff txs,db
	addtxs := []*models.EosdaqTx{}
	addvol := uint(0)
	for _, t := range txs {
		if _, ok := txMaps[t.ID]; !ok {
			addtxs = append(addtxs, t)
			addvol += t.GetVolume(eu.ticker.TokenSymbol)
		}
	}

	if len(addtxs) == 0 {
		return nil
	}
	//mlog.Debugw("UpdateTransactions db add", "cont", eu.ticker.ContractAccount, "data", addtxs)

	if err = eu.eosdaqRepo.SaveTransaction(innerCtx, addtxs); err != nil {
		mlog.Errorw("UpdateTransaction", "contract", eu.ticker.ContractAccount, "txs", addtxs, "err", err)
		return err
	}

	eu.ticker.CurrentPrice = addtxs[len(addtxs)-1].Price
	eu.ticker.Volume += addvol
	if err = eu.eosdaqRepo.UpdateTicker(innerCtx, eu.ticker); err != nil {
		mlog.Errorw("UpdateTicker", "contract", eu.ticker.ContractAccount, "ticker", eu.ticker, "err", err)
		return err
	}

	return nil
}
