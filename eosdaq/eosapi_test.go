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
		{"normal", args{&eosNet{"http://10.100.100.2", 18888, "cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"}, []string{"5HtZU5SArLEK3WDNntrK9fRCU8GFm9Ga4EAt9omGuYwiiFxMRyd"}}, false},
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
