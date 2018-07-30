package eosdaq

import (
	"fmt"
	"time"
)

type eosTimer struct {
	ticker  *time.Ticker
	eventCh chan struct{}
}

var eosTimer eosTimer

func init() {
	eosTimer.TimerOn(time.Millisecond * 500)
}

func (e eosTimer) TimerOn(d time.Duration) {
	e.TimerOff()
	e.ticker = time.NewTicker(d)
	e.eventCh = makeEventHandler()
	go eosTimer(e.ticker, e.eventCh)
}

func (e eosTimer) TimerOff() {
	e.ticker.Stop()
	close(e.eventCh)
}

func eosTimer(ticker *time.Ticker, ch chan<- struct{}) {
	for t := range ticker.C {
		ch <- struct{}{}
	}
}

func makeEventHandler() (ch chan struct{}) {
	ch = make(chan struct{})
	go func() {
		for {
			select {
			case <-ch:
				fmt.Println("event!")
			}
		}
	}()
	return ch
}
