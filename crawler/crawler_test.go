package crawler

import (
	"testing"

	eos "github.com/eoscanada/eos-go"
)

func TestNewCrawler(t *testing.T) {
	type args struct {
		api *eos.API
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"normal", args{nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCrawler(tt.args.api)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCrawler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//got.Wakeup()
			got.Stop()
		})
	}
}
