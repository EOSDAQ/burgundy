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

	keyBag := eos.NewKeyBag()
	keys := strings.Split(burgundy.GetString(fmt.Sprintf("%s_key", eosnet.contract)), ",")
	for _, key := range keys {
		if err := keyBag.Add(key); err != nil {
			return nil, errors.Annotatef(err, "contract[%s] add key error [%s]", eosnet.contract, key)
		}
	}
	api.SetSigner(keyBag)
	return &EosdaqAPI{api, eosnet.contract}, nil
}

func (e *EosdaqAPI) DoAction(action *eos.Action) error {
	e.Debug = true
	resp, err := e.SignPushActions(action)
	e.Debug = false
	if err != nil {
		mlog.Infow("ERROR calling : ", "err", err)
	} else {
		mlog.Infow("RESP : ", "resp", resp)
	}
	return err
}

func (e *EosdaqAPI) GetTx() (result []*models.EosdaqTx) {
	var res []*models.EosdaqTx
	out := &eos.GetTableRowsResp{More: true}
	begin, end := uint(0), uint(0)
	for out.More {
		out, _ = e.GetTableRows(eos.GetTableRowsRequest{
			Scope: e.contract,
			Code:  e.contract,
			Table: "tx",
			JSON:  true,
		})
		if out == nil || len(res) == 0 {
			break
		}
		out.JSONToStructs(&res)
		result = append(result, res...)
		if begin == 0 {
			begin = res[0].ID
		}
		end = res[len(res)-1].ID
	}
	if len(result) > 0 {
		fmt.Printf("delete tx from[%d] to[%d]\n", begin, end)
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
	var res []*models.OrderBook
	out := &eos.GetTableRowsResp{More: true}
	for out.More {
		out, _ = e.GetTableRows(eos.GetTableRowsRequest{
			Scope: e.contract,
			Code:  e.contract,
			Table: orderType.String(),
			JSON:  true,
		})
		if out == nil || len(res) == 0 {
			break
		}
		out.JSONToStructs(&res)
		for _, r := range res {
			r.Type = orderType
		}
		result = append(result, res...)
	}
	return result
}
