package models

import (
	"reflect"
	"testing"
	"time"
)

func TestEosdaqTx_GetArgs(t *testing.T) {
	et := &EosdaqTx{
		ID:         1,
		Price:      100,
		Maker:      "m",
		MakerAsset: "abc",
		Taker:      "t",
		TakerAsset: "sys",
		OrderTime:  time.Now().UnixNano(),
	}
	tests := []struct {
		name string
		et   *EosdaqTx
		want []interface{}
	}{
		// TODO: Add test cases.
		{"normal", et, []interface{}{et.Price, et.Maker, et.MakerAsset, et.Taker, et.TakerAsset, et.OrderTime}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.et.GetArgs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EosdaqTx.GetArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
