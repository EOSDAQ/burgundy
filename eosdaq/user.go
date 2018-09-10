package eosdaq

import (
	eos "github.com/eoscanada/eos-go"
)

// RegisterAction ... push action eosdaq enroll '[ "eosdaq" ]' -p eosdaq@active
func (e *API) RegisterAction(account string) *eos.Action {
	return action(e.contract, e.manage, account, "enroll")
}

// UnregisterAction ...
func (e *API) UnregisterAction(account string) *eos.Action {
	return action(e.contract, e.manage, account, "drop")
}

func action(contract, manage, account, action string) *eos.Action {
	return &eos.Action{
		Account: AN(contract),
		Name:    ActN(action),
		Authorization: []eos.PermissionLevel{
			{Actor: AN(manage), Permission: PN("active")},
		},
		ActionData: eos.NewActionData(ActionData{
			Contract: AN(manage),
			Account:  eos.AccountName(account),
		}),
	}
}

// ActionData ...
type ActionData struct {
	Contract eos.AccountName `json:"owner"`
	Account  eos.AccountName `json:"name"`
}
