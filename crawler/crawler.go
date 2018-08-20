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
}

var mlog *zap.SugaredLogger

func init() {
	mlog, _ = util.InitLog("crawler", "console")
}

func getTokens(tokenRepo _Repo.TokenRepository) (tokens []*models.Token) {
	var err error
	tokens, err = tokenRepo.GetTokens(context.Background())
	if err != nil {
		mlog.Infow("getTokens error", "err", err)
		return nil
	}
	return tokens
}

func InitModule(burgundy conf.ViperConfig, cancel <-chan os.Signal, db *gorm.DB) error {

	host := burgundy.GetString("eos_host")
	port := burgundy.GetInt("eos_port")
	tokenRepo := _Repo.NewGormTokenRepository(burgundy, db)
	tokens := getTokens(tokenRepo)

	timeout := time.Duration(burgundy.GetInt("timeout")) * time.Second
	crawlTimer := time.Duration(burgundy.GetInt("eos_crawlMS")) * time.Millisecond

	for _, t := range tokens {
		eosnet := eosdaq.NewEosnet(host, port, t.ContractAccount)
		api, err := eosdaq.NewAPI(burgundy, eosnet)
		if err != nil {
			mlog.Infow("InitModule error", "token", t, "err", err)
			continue
			//return errors.Annotatef(err, "InitModule NewAPI failed token[%+v]", t)
		}

		eosRepo := _Repo.NewGormEosdaqRepository(burgundy, db, t.ContractAccount)
		eossvc, err := service.NewEosdaqService(burgundy, t, eosRepo, tokenRepo, timeout)
		if err != nil {
			return errors.Annotatef(err, "InitModule NewSvc failed token[%s]", t)
		}

		err = NewCrawler(api, eossvc, crawlTimer, cancel)
		if err != nil {
			return errors.Annotatef(err, "InitModule NewCrawler failed token[%s]", t)
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
