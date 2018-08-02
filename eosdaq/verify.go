package eosdaq

import (
	eos "github.com/eoscanada/eos-go"
)

// SyncVerify ... push action eosdaq verify '[ "eosdaq" ]' -p eosdaq@active
func SyncVerify(contract eos.AccountName) *eos.Action {
	return &eos.Action{
		Account: contract,
		Name:    ActN("validate"),
		Authorization: []eos.PermissionLevel{
			{Actor: contract, Permission: PN("active")},
		},
		ActionData: eos.NewActionData(Verify{
			Contract: contract,
		}),
	}
}

type Verify struct {
	Contract eos.AccountName `json:"contract"`
}
