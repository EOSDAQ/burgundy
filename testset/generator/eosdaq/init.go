package eosdaq

import (
	"github.com/eoscanada/eos-go"
)

func init() {
	eos.RegisterAction(AN("eosio"), ActN("verify"), Verify{})
}

var AN = eos.AN
var PN = eos.PN
var ActN = eos.ActN
