package eosdaq

import (
	"os"
	"time"
)

type eosTimer struct {
	ticker  *time.Ticker
	crawler *Crawler
}

func NewTimer(c *Crawler, d time.Duration, cancel <-chan os.Signal) (*eosTimer, error) {
	timer := &eosTimer{}
	timer.ticker = time.NewTicker(d)
	timer.crawler = c
	go eosTicker(timer, cancel)
	return timer, nil
}

func (e *eosTimer) TimerOff() {
	if e.ticker != nil {
		e.ticker.Stop()
		e.ticker = nil
	}
	if e.crawler != nil {
		e.crawler.Stop()
		e.crawler = nil
	}
}

func eosTicker(et *eosTimer, cancel <-chan os.Signal) {
	for {
		select {
		case <-et.ticker.C:
			et.crawler.Wakeup()
		case <-cancel:
			et.TimerOff()
			break
		}
	}
}
