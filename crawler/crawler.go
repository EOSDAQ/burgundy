package crawler

import (
	"burgundy/conf"
	"burgundy/eosdaq"
	_Repo "burgundy/repository"
	"burgundy/service"
	"context"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
)

type Crawler struct {
	api           *eosdaq.EosdaqAPI
	EosdaqService service.EosdaqService
	ticker        *time.Ticker
}

func InitModule(burgundy conf.ViperConfig, cancel <-chan os.Signal, db *gorm.DB) error {

	host := burgundy.GetString("eos_host")
	port := burgundy.GetInt("eos_port")
	contracts := strings.Split(burgundy.GetString("eos_contract"), ",")

	timeout := time.Duration(burgundy.GetInt("timeout")) * time.Second
	crawlTimer := time.Duration(burgundy.GetInt("eos_crawl")) * time.Millisecond

	for _, c := range contracts {
		eosnet := &eosdaq.EosNet{
			host:     host,
			port:     port,
			contract: c,
		}

		api, err := NewAPI(burgundy, eosnet)
		if err != nil {
			return errors.Annotatef(err, "InitModule NewAPI failed contract[%s]", c)
		}

		eosRepo := _Repo.NewGormEosdaqRepository(db, c)
		eossvc, err := service.NewEosdaqService(burgundy, eosRepo, timeout)
		if err != nil {
			return errors.Annotatef(err, "InitModule NewSvc failed contract[%s]", c)
		}

		err := NewCrawler(api, eossvc, crawlTimer, cancel)
		if err != nil {
			return errors.Annotatef(err, "InitModule NewCrawler failed contract[%s]", c)
		}
	}

	return nil
}

func NewCrawler(api *eosdaq.EosdaqAPI, eosdaq service.EosdaqService, d time.Duration, cancel <-chan os.Signal) error {
	c := &Crawler{
		api:           api,
		EosdaqService: eosdaq,
	}
	return c.runCrawler(d, cancel)
}

func (c *Crawler) runCrawler(d time.Duration, cancel <-chan os.Signal) error {
	go func(ic *Crawler, d time.Duration) {
		t := time.NewTicker(d)
		for _ = range t.C {
			ctx := context.Background()
			ic.EosdaqService.UpdateOrderbook(ctx, ic.api.GetAsk())
		}
	}(c, d)
	go func(ic *Crawler) {
		t := time.NewTicker(d)
		for _ = range t.C {
			ctx := context.Background()
			ic.EosdaqService.UpdateOrderbook(ctx, ic.api.GetBid())
		}
	}(c, d)
	go func(ic *Crawler) {
		t := time.NewTicker(d)
		for _ = range t.C {
			ctx := context.Background()
			ic.EosdaqService.UpdateTransaction(ctx, ic.api.GetTx())
		}
	}(c, d)
	return nil
}
