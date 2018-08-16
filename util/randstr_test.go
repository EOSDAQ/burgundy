package util

import (
	"testing"
)

func TestConvertBase(t *testing.T) {
	type args struct {
		n    int
		base int
	}
	tests := []struct {
		name       string
		args       args
		wantString string
	}{
		// TODO: Add test cases.
		{"3 base", args{100, 3}, "10201"},
		{"2 base", args{16, 2}, "10000"},
		{"6 base", args{0, 6}, "0"},
		{"6 base", args{6, 6}, "10"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertBase(tt.args.n, tt.args.base); got != tt.wantString {
				t.Errorf("ConvertBase() = %v, want %v", got, tt.wantString)
			}
		})
	}
}
