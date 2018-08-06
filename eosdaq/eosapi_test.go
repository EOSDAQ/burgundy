package eosdaq

import (
	"testing"
)

func TestNewAPI(t *testing.T) {
	type args struct {
		eosnet *eosNet
		keys   []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"normal", args{&eosNet{"http://eosnet", 18888, "eosdaq", "eosdaqacnt"}, []string{"5KF7GfYH3SXjo4MNUUG8qg3Cs72T9mZueM12cvUTgBv8621VXcL"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAPI(tt.args.eosnet, tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("NewAPI() = %v", got)
		})
	}
}
