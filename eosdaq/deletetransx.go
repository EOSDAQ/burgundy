package eosdaq

import (
	eos "github.com/eoscanada/eos-go"
)

type Transx struct {
	Contract eos.AccountName `json:"name"`
	Begin    uint64          `json:"baseId"`
	End      uint64          `json:"endId"`
}

// DeleteTransaction ... push action eosdaq deletetransx '[ "eosdaq", 0, 0 ]' -p eosdaq@active
func DeleteTransaction(contract, manage eos.AccountName, begin, end uint) *eos.Action {
	return &eos.Action{
		Account: contract,
		Name:    ActN("deletetransx"),
		Authorization: []eos.PermissionLevel{
			{Actor: manage, Permission: PN("active")},
		},
		ActionData: eos.NewActionData(Transx{ //`["eosdaq",0,0]`),
			Contract: manage,
			Begin:    uint64(begin),
			End:      uint64(end),
		}),
	}
}
