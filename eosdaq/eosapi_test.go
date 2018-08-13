package eosdaq

import (
	"testing"

	"burgundy/conf"
	eos "github.com/eoscanada/eos-go"
)

Burgundy := conf.Burgundy

func TestEOSAPI(t *testing.T) {

	eosapi, err := NewAPI(Burgundy, &EosNet{
		host: "http://10.100.100.2",
		port: 18888,
		contract: "eosdaq",
	})
	if err != nil {
		t.Errorf("NewAPI failed [%s]", err)
	}

	eosapi.DoAction(RegisterAction("eosdaqacnt","newrovp"))
	eosapi.DoAction(UnregisterAction("eosdaqacnt","newrovp"))
	_ = eosapi.GetTx()
	_ = eosapi.GetAsk()
	_ = eosapi.GetBid()
}
