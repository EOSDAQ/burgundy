package eosdaq

import (
	"testing"
	"time"
)

func Test_eosTimer_OnOff(t *testing.T) {
	receiver := make(chan struct{})
	timer, _ := NewTimer(&Crawler{receiver, nil}, time.Millisecond*1, nil)
	select {
	case <-receiver:
		t.Log("OK")
	case <-time.After(time.Second * 2):
		t.Log("NOK")
	}
	timer.TimerOff()
}
