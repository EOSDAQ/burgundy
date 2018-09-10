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
	token      *models.Token
	eosdaqRepo repository.EosdaqRepository
	tokenRepo  repository.TokenRepository
	ctxTimeout time.Duration
}

// NewEosdaqService ...
func NewEosdaqService(burgundy *conf.ViperConfig,
	t *models.Token,
	er repository.EosdaqRepository,
	tr repository.TokenRepository,
	timeout time.Duration) (EosdaqService, error) {
	return &eosdaqUsecase{
		token:      t,
		eosdaqRepo: er,
		tokenRepo:  tr,
		ctxTimeout: timeout,
	}, nil
}

// UpdateOrderbook ...
func (eu eosdaqUsecase) UpdateOrderbook(ctx context.Context, obs []*models.OrderBook, orderType models.OrderType) (err error) {

	innerCtx, cancel := context.WithTimeout(ctx, eu.ctxTimeout)
	defer cancel()

	// get db old
	orderBooks, err := eu.eosdaqRepo.GetOrderBook(innerCtx, orderType)
	if err != nil {
		mlog.Errorw("UpdateOrderbook get", "contract", eu.token.ContractAccount, "err", err)
		return err
	}
	//mlog.Debugw("UpdateOrderbook db read", "cont", eu.token.ContractAccount, "data", orderBooks)
	orderMaps := make(map[string]*models.OrderBook)
	for _, o := range orderBooks {
		orderMaps[fmt.Sprintf("%d.%s", o.ID, o.OrderSymbol)] = o
	}
	// diff obs,db
	addBooks := []*models.OrderBook{}
	updBooks := []*models.OrderBook{}
	for _, n := range obs {
		key := fmt.Sprintf("%d.%s", n.ID, n.OrderSymbol)
		if o, ok := orderMaps[key]; ok {
			if o.Volume == n.Volume {
				delete(orderMaps, key)
			} else {
				updBooks = append(updBooks, n)
			}
		} else {
			addBooks = append(addBooks, n)
		}
	}
	//mlog.Debugw("UpdateOrderbook db add", "cont", eu.token.ContractAccount, "data", addBooks)

	// insert collection
	if err = eu.eosdaqRepo.SaveOrderBook(innerCtx, addBooks); err != nil {
		mlog.Errorw("UpdateOrderbook save", "contract", eu.token.ContractAccount, "err", err, "add", addBooks)
		return err
	}

	if err = eu.eosdaqRepo.UpdateOrderBook(innerCtx, updBooks); err != nil {
		mlog.Errorw("UpdateOrderbook save", "contract", eu.token.ContractAccount, "err", err, "add", addBooks)
		return err
	}

	delBooks := []*models.OrderBook{}
	for _, d := range orderMaps {
		delBooks = append(delBooks, d)
	}
	//mlog.Debugw("UpdateOrderbook db del", "cont", eu.token.ContractAccount, "data", delBooks)
	// delete collection
	if err = eu.eosdaqRepo.DeleteOrderBook(innerCtx, delBooks); err != nil {
		mlog.Errorw("UpdateOrderbook delete", "contract", eu.token.ContractAccount, "err", err, "del", delBooks)
		return err
	}

	// websocket broadcast

	return
}

// GetLastTransaction ...
func (eu eosdaqUsecase) GetLastTransactionID(ctx context.Context) (lastIdx int64) {
	innerCtx, cancel := context.WithTimeout(ctx, eu.ctxTimeout)
	defer cancel()

	return eu.eosdaqRepo.GetLastTransactionID(innerCtx)
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
		mlog.Errorw("UpdateTransactions get", "contract", eu.token.ContractAccount, "err", err)
		return err
	}
	//mlog.Debugw("UpdateTransaction db read", "cont", eu.token.ContractAccount, "data", dbtxs)
	txMaps := make(map[int64]struct{})
	for _, t := range dbtxs {
		txMaps[t.ID] = struct{}{}
	}

	// diff txs,db
	addtxs := []*models.EosdaqTx{}
	addvol := uint64(0)
	for _, t := range txs {
		if _, ok := txMaps[t.ID]; !ok {
			addtxs = append(addtxs, t)
			addvol += t.GetVolume(eu.token.Symbol)
		}
	}

	if len(addtxs) == 0 {
		return nil
	}
	//mlog.Debugw("UpdateTransactions db add", "cont", eu.token.ContractAccount, "data", addtxs)

	if err = eu.eosdaqRepo.SaveTransaction(innerCtx, addtxs); err != nil {
		mlog.Errorw("UpdateTransaction", "contract", eu.token.ContractAccount, "txs", addtxs, "err", err)
		return err
	}

	eu.token.CurrentPrice = addtxs[len(addtxs)-1].Price
	eu.token.Volume += addvol
	if err = eu.tokenRepo.UpdateToken(innerCtx, eu.token); err != nil {
		mlog.Errorw("UpdateToken", "contract", eu.token.ContractAccount, "token", eu.token, "err", err)
		return err
	}

	return nil
}
