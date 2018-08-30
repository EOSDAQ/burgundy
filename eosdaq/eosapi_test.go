package eosdaq

import (
	"testing"

	"burgundy/conf"

	"github.com/stretchr/testify/assert"
)

func TestEOSAPI(t *testing.T) {

	Burgundy := conf.Burgundy
	Burgundy.Set("eosdaqmanage", "5KX4aFJCuqyndWJnBLNanBxLffrT3ACeufxJL1N931V9uAU1Nnm")

	eosapi, err := NewAPI(Burgundy, &EosNet{
		host:     "http://10.100.100.2",
		port:     18888,
		contract: "eosdaq555555",
		manage:   "eosdaqmanage",
	})
	eosapi.Debug = false
	if err != nil {
		t.Errorf("NewAPI failed [%s]", err)
	}

	eosapi.DoAction(eosapi.RegisterAction("newrovp"))
	eosapi.DoAction(eosapi.UnregisterAction("newrovp"))

	eosapi.contract = "eosdaqoooo2o"

	assert.Nil(t, eosapi.GetTx(0))
	assert.NotNil(t, eosapi.GetAsk())
	assert.NotNil(t, eosapi.GetBid())
	assert.NotNil(t, eosapi.GetActions(0))
}
