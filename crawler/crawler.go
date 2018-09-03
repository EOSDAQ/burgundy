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
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"go.uber.org/zap"
)

type Crawler struct {
	api           *eosdaq.EosdaqAPI
	EosdaqService service.EosdaqService
	token         string
}

type crawlerDataHandler struct {
	begin int64
	end   int64
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

func InitModule(burgundy *conf.ViperConfig, cancel <-chan os.Signal, db *gorm.DB) error {

	host := burgundy.GetString("eos_host")
	port := burgundy.GetInt("eos_port")
	tokenRepo := _Repo.NewGormTokenRepository(burgundy, db)
	tokens := getTokens(tokenRepo)

	crawlContract := strings.Split(burgundy.GetString("eos_crawlexclude"), ",")
	crawlMap := make(map[string]struct{})
	for _, c := range crawlContract {
		if c == "" {
			continue
		}
		crawlMap[c] = struct{}{}
	}

	timeout := time.Duration(burgundy.GetInt("timeout")) * time.Second
	crawlTimer := time.Duration(burgundy.GetInt("eos_crawlMS")) * time.Millisecond
	manageContract := burgundy.GetString("eos_managecontract")

	for _, t := range tokens {
		if _, ok := crawlMap[t.ContractAccount]; ok {
			continue
		}
		eosnet := eosdaq.NewEosnet(host, port, t.ContractAccount, manageContract)
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

		err = NewCrawler(api, eossvc, t.Symbol, crawlTimer, cancel)
		if err != nil {
			return errors.Annotatef(err, "InitModule NewCrawler failed token[%s]", t)
		}
	}

	return nil
}

func NewCrawler(api *eosdaq.EosdaqAPI, eosdaq service.EosdaqService, token string,
	d time.Duration, cancel <-chan os.Signal) error {
	c := &Crawler{
		api:           api,
		EosdaqService: eosdaq,
		token:         token,
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
		cdh := &crawlerDataHandler{int64(1), int64(0)}
		for _ = range t.C {
			ctx := context.Background()
			ic.EosdaqService.UpdateTransaction(ctx, cdh.GetRangeData(ic.api, ic.token))
		}
	}(c, d)
	return nil
}

func (cdh *crawlerDataHandler) GetRangeData(api *eosdaq.EosdaqAPI, token string) (result []*models.EosdaqTx) {
	result = api.GetActionTxs(cdh.end, token)
	if len(result) == 0 {
		return nil
	}

	cdh.end = result[len(result)-1].ID
	/*
		if cdh.end-cdh.begin+1 >= 100 {
			mlog.Infow("delete tx", "from", cdh.begin, "to", cdh.end)
			api.DelTx(cdh.begin, cdh.end-1)
			cdh.begin = cdh.end
		}
	*/

	return result
}
