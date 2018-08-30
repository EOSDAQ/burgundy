package models

import (
	"fmt"
	"testing"

	eos "github.com/eoscanada/eos-go"
)

func TestEosdaqTx_GetVolume(t *testing.T) {
	type args struct {
		symbol string
	}
	tests := []struct {
		name string
		et   EosdaqTx
		args args
		want uint
	}{
		{"same symbol", EosdaqTx{ParseData: ParseData{Volume: uint(123123), Symbol: "ABC"}}, args{"ABC"}, uint(123123)},
		{"diff symbol", EosdaqTx{ParseData: ParseData{Volume: uint(123123), Symbol: "ABC"}}, args{"DEF"}, uint(0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.et.GetVolume(tt.args.symbol); got != tt.want {
				t.Errorf("EosdaqTx.GetVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContractData_Parse(t *testing.T) {
	tests := []struct {
		name string
		args ContractData
		want ParseData
	}{
		{"err", ContractData{"newro", "eosdaq", "123123 SYS", "0.0030"}, nil},
		{"bid", ContractData{"newro", "eosdaq", "123.123 SYS", "0.0030"}, ParseData{"newro", 123123, "SYS", BID, 30}},
		{"ask", ContractData{"newro", "eosdaq", "123.123 IPOS", "0.0030"}, ParseData{"newro", 123123, "IPOS", ASK, 30}},
		{"system", ContractData{"eosio.ram", "eosdaq", "123.123 SYS", "buy ram"}, nil},
		{"system", ContractData{"eosdaq", "eosio.ramfee", "123.123 SYS", "buy ram"}, nil},
		{"system", ContractData{"eosdaq", "eosio.msig", "123.123 SYS", "buy ram"}, nil},
		{"system", ContractData{"eosio.stake", "eosdaq", "123.123 SYS", "buy ram"}, nil},
		{"system", ContractData{"eosdaq", "eosio.token", "123.123 SYS", "buy ram"}, nil},
		{"system", ContractData{"eosio.saving", "eosdaq", "123.123 SYS", "buy ram"}, nil},
		{"system", ContractData{"eosdaq", "eosio.names", "123.123 SYS", "buy ram"}, nil},
		{"system", ContractData{"eosio.bpay", "eosdaq", "123.123 SYS", "buy ram"}, nil},
		{"system", ContractData{"eosio.vpay", "eosdaq", "123.123 SYS", "buy ram"}, nil},
		{"match1", ContractData{"eosdaq", "newroask", "123.123 SYS", "matched@0.0030"}, ParseData{"newroask", 123123, "SYS", MATCH, 30}},
		{"match2", ContractData{"eosdaq", "newrobid", "123.123 IPOS", "matched@0.0030"}, ParseData{"newrobid", 123123, "IPOS", MATCH, 30}},
		{"cancel1", ContractData{"eosdaq", "newrobid", "123.123 SYS", "cancel@0.0030"}, ParseData{"newrobid", 123123, "SYS", CANCEL, 30}},
		{"cancel2", ContractData{"eosdaq", "newroask", "123.123 IPOS", "cancel@0.0030"}, ParseData{"newroask", 123123, "IPOS", CANCEL, 30}},
		{"refund1", ContractData{"eosdaq", "newrobid", "123.123 SYS", "refund"}, ParseData{"newrobid", 123123, "SYS", REFUND, 30}},
		{"refund2", ContractData{"eosdaq", "newroask", "123.123 IPOS", "refund"}, ParseData{"newroask", 123123, "IPOS", REFUND, 30}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.Parse("eosdaq"); got != tt.want {
				t.Errorf("ContractData.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ContractData(t *testing.T) {
	cd := ContractData{
		From:     "newro",
		To:       "contract",
		Quantity: "123.0123 SYS",
		Memo:     "matched@0.0040",
	}
	ad := eos.NewActionData(cd)
	recd := ad.Data.(ContractData)
	fmt.Println(recd)
}
