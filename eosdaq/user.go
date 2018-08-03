package eosdaq

import (
	eos "github.com/eoscanada/eos-go"
)

// RegisterAction ... push action eosdaq enroll '[ "eosdaq" ]' -p eosdaq@active
func RegisterAction(contract, account string) *eos.Action {
	return action(contract, account, "enroll")
}

func UnregisterAction(contract, account string) *eos.Action {
	return action(contract, account, "drop")
}

func action(contract, account, action string) *eos.Action {
	eContract := eos.AccountName(contract)
	return &eos.Action{
		Account: eContract,
		Name:    ActN(action),
		Authorization: []eos.PermissionLevel{
			{Actor: eContract, Permission: PN("active")},
		},
		ActionData: eos.NewActionData(EosdaqAction{
			Contract: eContract,
			Account:  eos.AccountName(account),
		}),
	}
}

type EosdaqAction struct {
	Contract eos.AccountName `json:"owner"`
	Account  eos.AccountName `json:"name"`
}
