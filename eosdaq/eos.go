package eosdaq

import (
	"encoding/hex"
	"fmt"
	"time"

	eos "github.com/eoscanada/eos-go"
	"github.com/juju/errors"
)

var AN = eos.AN
var PN = eos.PN
var ActN = eos.ActN

func init() {
	eos.RegisterAction(AN("eosio"), ActN("verify"), Verify{})
}

type eosNet struct {
	host    string
	port    int
	chainID string
}

/*
   "id": 0,
   "price": 30,
   "maker": "newrotaker",
   "maker_asset": "0.0011 SYS",
   "taker": "newrovp",
   "taker_asset": "0.0333 ABC",
   "ordertime": 739407904
*/
type EosdaqTx struct {
	ID         int       `json:"id"`
	Price      int       `json:"price"`
	Maker      string    `json:"maker"`
	MakerAsset string    `json:"maker_asset"`
	Taker      string    `json:"taker"`
	TakerAsset string    `json:"taker_asset"`
	OrderTime  time.Time `json:"ordertime"`
}

func NewAPI(eosnet *eosNet, keys []string) (*eos.API, error) {
	cid, err := hex.DecodeString(eosnet.chainID)
	if err != nil {
		return nil, errors.Annotatef(err, "invalid chain ID[%s]", eosnet.chainID)
	}
	api := eos.New(fmt.Sprintf("%s:%d", eosnet.host, eosnet.port), cid)
	keyBag := eos.NewKeyBag()
	for _, key := range keys {
		if err = keyBag.Add(key); err != nil {
			return nil, errors.Annotatef(err, "add key error [%s]", key)
		}
	}
	api.SetSigner(keyBag)
	return api, nil
}
