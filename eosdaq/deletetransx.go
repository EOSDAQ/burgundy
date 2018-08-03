package eosdaq

import (
	eos "github.com/eoscanada/eos-go"
)

// DeleteTransaction ... push action eosdaq deletetransx '[ "eosdaq", 0, 0 ]' -p eosdaq@active
func DeleteTransaction(contract eos.AccountName, begin, end int) *eos.Action {
	return &eos.Action{
		Account: contract,
		Name:    ActN("deletetransx"),
		Authorization: []eos.PermissionLevel{
			{Actor: contract, Permission: PN("active")},
		},
		ActionData: eos.NewActionData(Transx{ //`["eosdaq",0,0]`),
			Contract: contract,
			Begin:    uint64(begin),
			End:      uint64(end),
		}),
	}
}

type Transx struct {
	Contract eos.AccountName `json:"name"`
	Begin    uint64          `json:"baseId"`
	End      uint64          `json:"endId"`
}
