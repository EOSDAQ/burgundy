package eosdaq

import (
	"burgundy/conf"
	"burgundy/models"
	"fmt"
	"strings"

	eos "github.com/eoscanada/eos-go"
	"github.com/juju/errors"
)

var AN = eos.AN
var PN = eos.PN
var ActN = eos.ActN

func init() {
	eos.RegisterAction(AN("eosio"), ActN("enroll"), EosdaqAction{})
	eos.RegisterAction(AN("eosio"), ActN("drop"), EosdaqAction{})
	eos.RegisterAction(AN("eosio"), ActN("deletetransx"), Transx{})
	//eos.Debug = true
}

type EosNet struct {
	host     string
	port     int
	contract string
	manage   string
}

func NewEosnet(host string, port int, contract, manage string) *EosNet {
	return &EosNet{host, port, contract, manage}
}

type EosdaqAPI struct {
	*eos.API
	contract string
	manage   string
}

func NewAPI(burgundy conf.ViperConfig, eosnet *EosNet) (*EosdaqAPI, error) {
	api := eos.New(fmt.Sprintf("%s:%d", eosnet.host, eosnet.port))

	keys := strings.Split(burgundy.GetString(eosnet.manage), ",")
	if keys[0] == "" {
		return nil, errors.Errorf("NewAPI no keys", "contract", eosnet.manage)
	} else {
		keyBag := eos.NewKeyBag()
		for _, key := range keys {
			if err := keyBag.Add(key); err != nil {
				return nil, errors.Annotatef(err, "New API contract[%s] add key error [%s]", eosnet.contract, key)
			}
		}
		api.SetSigner(keyBag)
	}

	if burgundy.GetString("loglevel") == "debug" {
		api.Debug = true
	}
	return &EosdaqAPI{api, eosnet.contract, eosnet.manage}, nil
}

func (e *EosdaqAPI) DoAction(action *eos.Action) error {
	resp, err := e.SignPushActions(action)
	if err != nil {
		mlog.Infow("ERROR calling : ", "err", err)
	} else {
		mlog.Infow("RESP : ", "resp", resp)
	}
	return err
}

func (e *EosdaqAPI) GetTx(start uint) (result []*models.EosdaqTx) {
	var err error
	out := &eos.GetTableRowsResp{More: true}
	end := start
	for out.More {
		out, err = e.GetTableRows(eos.GetTableRowsRequest{
			Scope:      e.contract,
			Code:       e.contract,
			LowerBound: fmt.Sprintf("%d", end+1),
			Table:      "tx",
			JSON:       true,
		})
		if err != nil {
			mlog.Errorw("GetTx error", "contract", e.contract, "err", err)
			break
		}
		if out == nil {
			mlog.Infow("GetTx nil", "contract", e.contract)
			break
		}
		res := models.TxResponse{}
		out.JSONToStructs(&res)
		if len(res) == 0 {
			//mlog.Infow("GetTx nil", "contract", e.contract)
			break
		}
		end = res[len(res)-1].ID
		//begin, end = res.GetRange(begin, end)
		//mlog.Infow("GetTx ", "b", begin, "e", end)
		result = append(result, res...)
	}
	/*
		if end-begin+1 > 100 {
			mlog.Infow("delete tx ", "from", begin, "to", end)
			e.DoAction(
				DeleteTransaction(eos.AccountName(e.contract), begin, end-1),
			)
		}
	*/
	return result
}

func (e *EosdaqAPI) DelTx(from, to uint) {
	e.DoAction(
		DeleteTransaction(eos.AccountName(e.contract), eos.AccountName(e.manage), from, to),
	)
}

func (e *EosdaqAPI) GetAsk() (result []*models.OrderBook) {
	return e.getOrderBook(models.ASK)
}
func (e *EosdaqAPI) GetBid() (result []*models.OrderBook) {
	return e.getOrderBook(models.BID)
}

func (e *EosdaqAPI) getOrderBook(orderType models.OrderType) (result []*models.OrderBook) {
	var err error
	out := &eos.GetTableRowsResp{More: true}
	begin, end := uint(0), uint(0)
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
		res := []*models.OrderBook{}
		out.JSONToStructs(&res)
		if len(res) == 0 {
			//mlog.Infow("getOrderBook nil", "contract", e.contract, "type", orderType)
			break
		}
		if begin == 0 {
			begin = res[0].ID
		}
		end = res[len(res)-1].ID
		for _, r := range res {
			r.Type = orderType
		}
		result = append(result, res...)
	}
	return result
}
