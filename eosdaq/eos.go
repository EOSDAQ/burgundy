package eosdaq

import (
	"fmt"
	"strconv"
	"time"

	eos "github.com/eoscanada/eos-go"
	"github.com/juju/errors"
)

var AN = eos.AN
var PN = eos.PN
var ActN = eos.ActN

func init() {
	eos.RegisterAction(AN("eosio"), ActN("validate"), Verify{})
	eos.RegisterAction(AN("eosio"), ActN("deletetransx"), Transx{})
	//eos.Debug = true
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
	OrderTime  Timestamp `json:"ordertime"`
}
type Timestamp struct {
	time.Time
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := t.Time.Unix()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}

	t.Time = time.Unix(int64(ts), 0)

	return nil
}

func NewAPI(eosnet *eosNet, keys []string) (*eos.API, error) {
	api := eos.New(fmt.Sprintf("%s:%d", eosnet.host, eosnet.port))
	infoResp, _ := api.GetInfo()
	mlog.Infow("NewAPI", "info", infoResp)
	accResp, _ := api.GetAccount("eosdaq")
	mlog.Infow("NewAPI", "acct", accResp)

	//wallet := eos.NewWalletSigner(api, "wall2")

	keyBag := eos.NewKeyBag()
	for _, key := range keys {
		if err := keyBag.Add(key); err != nil {
			//if err = wallet.ImportPrivateKey(key); err != nil {
			return nil, errors.Annotatef(err, "add key error [%s]", key)
		}
	}
	api.SetSigner(keyBag)
	//api.SetSigner(wallet)
	return api, nil
}
