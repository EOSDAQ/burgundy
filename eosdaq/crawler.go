package eosdaq

import (
	"encoding/json"
	"fmt"
	"time"

	eos "github.com/eoscanada/eos-go"
)

type Crawler struct {
	receiver chan struct{}
	api      *eos.API
}

func NewCrawler(api *eos.API) (*Crawler, error) {
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
	//var res []*EosdaqTx
	out := &eos.GetTableRowsResp{More: true}
	for out.More {
		out, _ = c.api.GetTableRows(eos.GetTableRowsRequest{
			Scope: "eosdaq",
			Code:  "eosdaq",
			Table: "tx",
			JSON:  true,
		})
		//out.BinaryToStructs(&res)
		//fmt.Printf("tx value [%v]\n", res)
		data, _ := json.Marshal(out)
		fmt.Printf("row [%s]\n", string(data))
	}
}
