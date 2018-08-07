package eosdaq

import (
	"burgundy/models"
	"fmt"

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

type EosdaqAPI struct {
	*eos.API
	eoscontract  string
	acctcontract string
}

func NewAPI(eosnet *EosNet, keys []string) (*EosdaqAPI, error) {
	api := eos.New(fmt.Sprintf("%s:%d", eosnet.host, eosnet.port))

	/*
		infoResp, _ := api.GetInfo()
		mlog.Infow("NewAPI", "info", infoResp)
		accResp, _ := api.GetAccount("eosdaq")
		mlog.Infow("NewAPI", "acct", accResp)
	*/

	keyBag := eos.NewKeyBag()
	for _, key := range keys {
		if err := keyBag.Add(key); err != nil {
			return nil, errors.Annotatef(err, "add key error [%s]", key)
		}
	}
	api.SetSigner(keyBag)
	return &EosdaqAPI{api, eosnet.contract, eosnet.acctcontract}, nil
}

func (e *EosdaqAPI) CrawlData() {
	var res []models.EosdaqTx
	out := &eos.GetTableRowsResp{More: true}
	for out.More {
		out, _ = e.GetTableRows(eos.GetTableRowsRequest{
			Scope: e.eoscontract,
			Code:  e.eoscontract,
			Table: "tx",
			JSON:  true,
		})
		if out == nil {
			break
		}
		out.JSONToStructs(&res)
		for _, r := range res {
			fmt.Printf("tx value [%v]\n", r)
		}
		if len(res) > 0 {
			begin, end := res[0].ID, res[len(res)-1].ID
			fmt.Printf("delete tx from[%d] to[%d]\n", begin, end)
			e.call(
				DeleteTransaction(eos.AccountName("eosdaq"), begin, end),
			)
		}
	}
}

func (e *EosdaqAPI) RegisterUser(account string) error {
	return e.call(RegisterAction(e.acctcontract, account))
}

func (e *EosdaqAPI) UnregisterUser(account string) error {
	return e.call(UnregisterAction(e.acctcontract, account))
}

func (e *EosdaqAPI) call(action *eos.Action) error {
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
