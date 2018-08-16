package models

import (
	"reflect"
	"testing"
)

func TestEosdaqTx_GetArgs(t *testing.T) {
	et := &EosdaqTx{
		ID:         1,
		Price:      100,
		Maker:      "m",
		MakerAsset: "abc",
		Taker:      "t",
		TakerAsset: "sys",
		OrderTime:  "12345678",
	}
	tests := []struct {
		name string
		et   *EosdaqTx
		want []interface{}
	}{
		// TODO: Add test cases.
		{"normal", et, []interface{}{et.ID, et.Price, et.Maker, et.MakerAsset, et.Taker, et.TakerAsset, et.OrderTime}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.et.GetArgs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EosdaqTx.GetArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		{"same symbol", EosdaqTx{MakerAsset: "123.123 ABC"}, args{"ABC"}, 1231230},
		{"same symbol", EosdaqTx{MakerAsset: "123.0000 SYS", TakerAsset: "123.123 ABC"}, args{"ABC"}, 1231230},
		{"diff symbol", EosdaqTx{TakerAsset: "123.123 ABC"}, args{"DEF"}, 0},
		{"similar symbol", EosdaqTx{MakerAsset: "123.123 SYS", TakerAsset: "123.123 ABC"}, args{"BCD"}, 0},
		{"bigger symbol", EosdaqTx{MakerAsset: "123.123 SYSABC", TakerAsset: "123.123 ABC"}, args{"SAB"}, 0},
		{"small symbol", EosdaqTx{MakerAsset: "123.123 IQ", TakerAsset: "123.123 ABC"}, args{"AIQ"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.et.GetVolume(tt.args.symbol); got != tt.want {
				t.Errorf("EosdaqTx.GetVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}
