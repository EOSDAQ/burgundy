package eosdaq

import (
	"burgundy/conf"
	"burgundy/models"
	"fmt"
	"strings"

	eos "github.com/eoscanada/eos-go"
	"github.com/juju/errors"
)

// EOS defines
var (
	AN   = eos.AN
	PN   = eos.PN
	ActN = eos.ActN
)

func init() {
	eos.RegisterAction(AN("eosio"), ActN("enroll"), ActionData{})
	eos.RegisterAction(AN("eosio"), ActN("drop"), ActionData{})
	//eos.Debug = true
}

// EosNet ...
type EosNet struct {
	host     string
	port     int
	contract string
	manage   string
}

// NewEosnet ...
func NewEosnet(host string, port int, contract, manage string) *EosNet {
	return &EosNet{host, port, contract, manage}
}

// API ...
type API struct {
	*eos.API
	contract string
	manage   string
}

// NewAPI ...
func NewAPI(burgundy *conf.ViperConfig, eosnet *EosNet) (*API, error) {
	api := eos.New(fmt.Sprintf("%s:%d", eosnet.host, eosnet.port))

	keys := strings.Split(burgundy.GetString(eosnet.manage), ",")
	if keys[0] == "" {
		return nil, errors.Errorf("NewAPI no keys contract[%s]", eosnet.manage)
	}

	keyBag := eos.NewKeyBag()
	for _, key := range keys {
		if err := keyBag.Add(key); err != nil {
			return nil, errors.Annotatef(err, "New API contract[%s] add key error [%s]", eosnet.contract, key)
		}
	}
	api.SetSigner(keyBag)

	if burgundy.GetString("loglevel") == "debug" {
		api.Debug = true
	}
	return &API{api, eosnet.contract, eosnet.manage}, nil
}

// DoAction ...
func (e *API) DoAction(action *eos.Action) error {
	resp, err := e.SignPushActions(action)
	if err != nil {
		mlog.Infow("ERROR calling : ", "err", err)
	} else {
		mlog.Infow("RESP : ", "resp", resp)
	}
	return err
}

// GetActionTxs ...
func (e *API) GetActionTxs(start int64, symbol string) (result []*models.EosdaqTx, end int64) {
	end = start
	out, err := e.GetActions(eos.GetActionsRequest{
		AccountName: AN(e.contract),
		Pos:         start + 1,
		Offset:      int64(100),
	})
	if err != nil {
		mlog.Errorw("GetActions error", "contract", e.contract, "err", err)
		return nil, end
	}
	if out == nil {
		mlog.Debugw("GetActions nil", "contract", e.contract)
		return nil, end
	}
	for _, o := range out.Actions {
		res := &models.EosdaqTx{
			ID:            o.AccountSeq,
			OrderSymbol:   symbol,
			OrderTime:     o.BlockTime.Time,
			TransactionID: o.Trace.TransactionID,
		}
		if end < res.ID {
			end = res.ID
		}

		if o.Trace.Action.ActionData.Data == nil {
			mlog.Debugw("GetActions nil data", "action", res)
			continue
		}

		cd := &models.ContractData{}
		res.EOSData = cd.MarshalData(o.Trace.Action.ActionData.Data)
		if res.EOSData == nil {
			mlog.Debugw("GetActions nil data", "action", o.Trace.Action.ActionData)
			continue
		}
		result = append(result, res)
	}
	return result, end
}

// GetAsk ...
func (e *API) GetAsk(symbol string) (result []*models.OrderBook) {
	return e.getOrderBook(symbol, models.ASK)
}

// GetBid ...
func (e *API) GetBid(symbol string) (result []*models.OrderBook) {
	return e.getOrderBook(symbol, models.BID)
}

func (e *API) getOrderBook(symbol string, orderType models.OrderType) (result []*models.OrderBook) {
	var err error
	out := &eos.GetTableRowsResp{More: true}
	end := uint(0)
	for out.More {
		out, err = e.GetTableRows(eos.GetTableRowsRequest{
			Scope:      e.contract,
			Code:       e.contract,
			LowerBound: fmt.Sprintf("%d", end+1),
			Table:      orderType.String(),
			JSON:       true,
		})
		if err != nil {
			mlog.Errorw("getOrderBook error", "contract", e.contract, "type", orderType, "err", err)
			break
		}
		if out == nil {
			//mlog.Infow("getOrderBook nil", "contract", e.contract, "type", orderType)
			break
		}
		res := []*models.OrderData{}
		out.JSONToStructs(&res)
		if len(res) == 0 {
			//mlog.Infow("getOrderBook nil", "contract", e.contract, "type", orderType)
			break
		}
		end = res[len(res)-1].ID
		for _, r := range res {
			ob := r.Parse(symbol, orderType)
			if ob != nil {
				result = append(result, ob)
			}
		}
	}
	return result
}
