package models

import (
	"reflect"
	"testing"
)

func TestEosdaqTx_GetVolume(t *testing.T) {
	type args struct {
		symbol string
	}
	tests := []struct {
		name string
		et   EosdaqTx
		args args
		want uint64
	}{
		{"same symbol", EosdaqTx{EOSData: &EOSData{Volume: uint64(123123), Symbol: "ABC"}}, args{"ABC"}, uint64(123123)},
		{"diff symbol", EosdaqTx{EOSData: &EOSData{Volume: uint64(123123), Symbol: "ABC"}}, args{"DEF"}, uint64(0)},
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
		cd   *ContractData
		want *EOSData
	}{
		{"err", &ContractData{"newro", "eosdaq", "123123 SYS", "0.0030"}, nil},
		//{"bid", &ContractData{"newro", "eosdaq", "123.1230 SYS", "0.0030"}, &EOSData{"newro", 1231230, "SYS", BID, 30}},
		//{"ask", &ContractData{"newro", "eosdaq", "123.1230 IPOS", "0.0030"}, &EOSData{"newro", 1231230, "IPOS", ASK, 30}},
		{"bid", &ContractData{"newro", "eosdaq", "123.1230 SYS", "0.0030"}, nil},
		{"ask", &ContractData{"newro", "eosdaq", "123.1230 IPOS", "0.0030"}, nil},
		{"system", &ContractData{"eosio.ram", "eosdaq", "123.1230 SYS", "buy ram"}, nil},
		{"system", &ContractData{"eosdaq", "eosio.ramfee", "123.1023 SYS", "buy ram"}, nil},
		{"system", &ContractData{"eosdaq", "eosio.msig", "123.1230 SYS", "buy ram"}, nil},
		{"system", &ContractData{"eosio.stake", "eosdaq", "123.1203 SYS", "buy ram"}, nil},
		{"system", &ContractData{"eosdaq", "eosio.token", "123.1203 SYS", "buy ram"}, nil},
		{"system", &ContractData{"eosio.saving", "eosdaq", "123.1023 SYS", "buy ram"}, nil},
		{"system", &ContractData{"eosdaq", "eosio.names", "123.1203 SYS", "buy ram"}, nil},
		{"system", &ContractData{"eosio.bpay", "eosdaq", "123.1230 SYS", "buy ram"}, nil},
		{"system", &ContractData{"eosio.vpay", "eosdaq", "123.1230 SYS", "buy ram"}, nil},
		{"match1", &ContractData{"eosdaq", "newroask", "123.1230 SYS", "match@0.0030"}, &EOSData{"newroask", 30, 1231230, "SYS", MATCH}},
		{"match2", &ContractData{"eosdaq", "newrobid", "123.1230 IPOS", "match@0.0030"}, &EOSData{"newrobid", 30, 1231230, "IPOS", MATCH}},
		//{"cancel1", &ContractData{"eosdaq", "newrobid", "123.1230 SYS", "cancel@0.0030"}, &EOSData{"newrobid", 1231230, "SYS", CANCEL, 30}},
		//{"cancel2", &ContractData{"eosdaq", "newroask", "123.1230 IPOS", "cancel@0.0030"}, &EOSData{"newroask", 1231230, "IPOS", CANCEL, 30}},
		{"cancel1", &ContractData{"eosdaq", "newrobid", "123.1230 SYS", "cancel@0.0030"}, nil},
		{"cancel2", &ContractData{"eosdaq", "newroask", "123.1230 IPOS", "cancel@0.0030"}, nil},
		//{"refund1", &ContractData{"eosdaq", "newrobid", "123.1230 SYS", "refund"}, &EOSData{"newrobid", 1231230, "SYS", REFUND, 0}},
		//{"refund2", &ContractData{"eosdaq", "newroask", "123.1230 IPOS", "refund"}, &EOSData{"newroask", 1231230, "IPOS", REFUND, 0}},
		{"refund1", &ContractData{"eosdaq", "newrobid", "123.1230 SYS", "refund"}, nil},
		{"refund2", &ContractData{"eosdaq", "newroask", "123.1230 IPOS", "refund"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cd.Parse("IPOS"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContractData.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
