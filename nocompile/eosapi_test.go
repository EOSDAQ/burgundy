package eosdaq

import (
	"testing"

	"burgundy/conf"

	"github.com/stretchr/testify/assert"
)

func _localTest(t *testing.T) {

	Burgundy := conf.Burgundy
	Burgundy.Set("eosdaqmanage", "5K8q9AzWV6ztfu16LHrngHG2Ts4SdzDPQhYCpTUC4Fx9jsnmBbo")

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

	t.Run("GetAsk", func(t *testing.T) {
		assert.NotNil(t, eosapi.GetAsk("IPOS"))
	})
	t.Run("GetBid", func(t *testing.T) {
		assert.NotNil(t, eosapi.GetBid("IPOS"))
	})
	t.Run("GetActionTxs", func(t *testing.T) {
		r, end := eosapi.GetActionTxs(int64(0), "IPOS")
		assert.NotNil(t, r)
		assert.EqualValues(t, r[len(r)-1].ID, end)
	})
}
