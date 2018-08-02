package eosdaq

import (
	eos "github.com/eoscanada/eos-go"
)

// RegisterAction ... push action eosdaq enroll '[ "eosdaq" ]' -p eosdaq@active
func RegisterAction(contract eos.AccountName) *eos.Action {
	return action(contract, "enroll")
}

func UnregisterAction(contract eos.AccountName) *eos.Action {
	return action(contract, "drop")
}

func action(contract eos.AccountName, action string) *eos.Action {
	return &eos.Action{
		Account: contract,
		Name:    ActN(action),
		Authorization: []eos.PermissionLevel{
			{Actor: contract, Permission: PN("active")},
		},
		ActionData: eos.NewActionData(EosdaqAction{
			Contract: contract,
		}),
	}
}

type EosdaqAction struct {
	Contract eos.AccountName `json:"contract"`
}
