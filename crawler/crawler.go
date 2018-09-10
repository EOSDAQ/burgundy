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

// Crawler ...
type Crawler struct {
	api           *eosdaq.API
	EosdaqService service.EosdaqService
	symbol        string
}

var mlog *zap.SugaredLogger

func init() {
	mlog, _ = util.InitLog("crawler", "devel")
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

// InitModule ...
func InitModule(burgundy *conf.ViperConfig, cancel <-chan os.Signal, db *gorm.DB) error {

	mlog, _ = util.InitLog("crawler", burgundy.GetString("loglevel"))

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

		eosRepo := _Repo.NewGormEosdaqRepository(db, t.Symbol)
		eossvc, err := service.NewEosdaqService(burgundy, t, eosRepo, tokenRepo, timeout)
		if err != nil {
			return errors.Annotatef(err, "InitModule NewSvc failed token[%s]", t.Symbol)
		}

		err = NewCrawler(api, eossvc, t.Symbol, crawlTimer, cancel)
		if err != nil {
			return errors.Annotatef(err, "InitModule NewCrawler failed token[%s]", t.Symbol)
		}
	}

	return nil
}

// NewCrawler ...
func NewCrawler(api *eosdaq.API, eosdaq service.EosdaqService, symbol string,
	d time.Duration, cancel <-chan os.Signal) error {
	c := &Crawler{
		api:           api,
		EosdaqService: eosdaq,
		symbol:        symbol,
	}
	return c.runCrawler(d, cancel)
}

func (c *Crawler) runCrawler(d time.Duration, cancel <-chan os.Signal) error {
	go func(ic *Crawler, d time.Duration) {
		t := time.NewTicker(d)
		for range t.C {
			ctx := context.Background()
			//mlog.Infow("Crawler UpdateOrderbook Bid")
			ic.EosdaqService.UpdateOrderbook(ctx, ic.api.GetBid(ic.symbol), models.BID)
		}
	}(c, d)
	go func(ic *Crawler, d time.Duration) {
		t := time.NewTicker(d)
		for range t.C {
			ctx := context.Background()
			//mlog.Infow("Crawler UpdateOrderbook Ask")
			ic.EosdaqService.UpdateOrderbook(ctx, ic.api.GetAsk(ic.symbol), models.ASK)
		}
	}(c, d)
	go func(ic *Crawler, d time.Duration) {
		t := time.NewTicker(d)
		lastIdx := ic.EosdaqService.GetLastTransactionID(context.Background())
		var result []*models.EosdaqTx
		for range t.C {
			ctx := context.Background()
			result, lastIdx = ic.api.GetActionTxs(lastIdx, ic.symbol)
			ic.EosdaqService.UpdateTransaction(ctx, result)
		}
	}(c, d)
	return nil
}
