package eosdaq

import (
	"fmt"
	"time"
)

type Crawler struct {
	receiver chan struct{}
	api      *EosdaqAPI
}

func NewCrawler(api *EosdaqAPI) (*Crawler, error) {
	c := &Crawler{
		api: api,
	}
	c.makeEventHandler()
	return c, nil
}

func (c *Crawler) makeEventHandler() {
	c.receiver = make(chan struct{})
	go func(innerCrawl *Crawler) {
		for {
			select {
			case <-innerCrawl.receiver:
				fmt.Println("event!", time.Now())
				innerCrawl.Do()
			}
		}
	}(c)
}

func (c *Crawler) Wakeup() {
	c.receiver <- struct{}{}
	//fmt.Println("Waketup", time.Now())
}

func (c *Crawler) Stop() {
	close(c.receiver)
}

func (c *Crawler) Do() {
	c.api.CrawlData()
}
