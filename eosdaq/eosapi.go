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
}

func NewEosnet(host string, port int, contract string) *EosNet {
	return &EosNet{host, port, contract}
}

type EosdaqAPI struct {
	*eos.API
	contract string
}

func NewAPI(burgundy conf.ViperConfig, eosnet *EosNet) (*EosdaqAPI, error) {
	api := eos.New(fmt.Sprintf("%s:%d", eosnet.host, eosnet.port))

	/*
		infoResp, _ := api.GetInfo()
		mlog.Infow("NewAPI", "info", infoResp)
		accResp, _ := api.GetAccount("eosdaq")
		mlog.Infow("NewAPI", "acct", accResp)
	*/

	keys := strings.Split(burgundy.GetString(eosnet.contract), ",")
	if keys[0] == "" {
		return nil, errors.Errorf("NewAPI no keys", "contract", eosnet.contract)
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
	return &EosdaqAPI{api, eosnet.contract}, nil
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

func (e *EosdaqAPI) GetTx() (result []*models.EosdaqTx) {
	var err error
	out := &eos.GetTableRowsResp{More: true}
	begin, end := uint(0), uint(0)
	for out.More {
		out, err = e.GetTableRows(eos.GetTableRowsRequest{
			Scope:      e.contract,
			Code:       e.contract,
			LowerBound: fmt.Sprintf("%d", end+2),
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
		begin, end = res.GetRange(begin, end)
		//mlog.Infow("GetTx ", "b", begin, "e", end)
		result = append(result, res...)
	}
	if len(result) > 100 {
		mlog.Infow("delete tx ", "from", begin, "to", end)
		e.DoAction(
			DeleteTransaction(eos.AccountName(e.contract), begin, end),
		)
	}
	return result
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
