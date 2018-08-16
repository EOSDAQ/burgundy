package crawler

import (
	"burgundy/conf"
	"burgundy/eosdaq"
	"burgundy/models"
	_Repo "burgundy/repository"
	"burgundy/service"
	"burgundy/util"
	"context"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"go.uber.org/zap"
)

type Crawler struct {
	api           *eosdaq.EosdaqAPI
	EosdaqService service.EosdaqService
	ticker        *time.Ticker
}

var mlog *zap.SugaredLogger

func init() {
	mlog, _ = util.InitLog("crawler", "console")
}

func getTickers(tickRepo _Repo.TickerRepository) (tickers []*models.Ticker) {
	var err error
	tickers, err = tickRepo.GetTickers(context.Background())
	if err != nil {
		mlog.Infow("getTickers error", "err", err)
		return nil
	}
	return tickers
}

func InitModule(burgundy conf.ViperConfig, cancel <-chan os.Signal, db *gorm.DB) error {

	host := burgundy.GetString("eos_host")
	port := burgundy.GetInt("eos_port")
	tickerRepo := _Repo.NewGormTickerRepository(burgundy, db)
	tickers := getTickers(tickerRepo)

	timeout := time.Duration(burgundy.GetInt("timeout")) * time.Second
	crawlTimer := time.Duration(burgundy.GetInt("eos_crawlMS")) * time.Millisecond

	for _, t := range tickers {
		eosnet := eosdaq.NewEosnet(host, port, t.ContractAccount)
		api, err := eosdaq.NewAPI(burgundy, eosnet)
		if err != nil {
			return errors.Annotatef(err, "InitModule NewAPI failed ticker[%+v]", t)
		}

		eosRepo := _Repo.NewGormEosdaqRepository(burgundy, db, t.ContractAccount)
		eossvc, err := service.NewEosdaqService(burgundy, t, eosRepo, tickerRepo, timeout)
		if err != nil {
			return errors.Annotatef(err, "InitModule NewSvc failed ticker[%s]", t)
		}

		err = NewCrawler(api, eossvc, crawlTimer, cancel)
		if err != nil {
			return errors.Annotatef(err, "InitModule NewCrawler failed ticker[%s]", t)
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
			//mlog.Infow("Crawler UpdateOrderbook Ask")
			ic.EosdaqService.UpdateOrderbook(ctx, ic.api.GetAsk(), models.ASK)
		}
	}(c, d)
	go func(ic *Crawler, d time.Duration) {
		t := time.NewTicker(d)
		for _ = range t.C {
			ctx := context.Background()
			//mlog.Infow("Crawler UpdateOrderbook Bid")
			ic.EosdaqService.UpdateOrderbook(ctx, ic.api.GetBid(), models.BID)
		}
	}(c, d)
	go func(ic *Crawler, d time.Duration) {
		t := time.NewTicker(d)
		for _ = range t.C {
			ctx := context.Background()
			ic.EosdaqService.UpdateTransaction(ctx, ic.api.GetTx())
		}
	}(c, d)
	return nil
}
