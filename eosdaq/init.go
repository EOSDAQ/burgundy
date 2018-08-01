package eosdaq

import (
	"burgundy/conf"
	"burgundy/util"
	"os"
	"time"

	"github.com/juju/errors"
	"go.uber.org/zap"
)

var mlog *zap.SugaredLogger
var timer *eosTimer

func init() {
	mlog, _ = util.InitLog("eosdaq", "console")
}

func InitModule(burgundy conf.ViperConfig, cancel <-chan os.Signal) error {

	eosnet := &eosNet{
		host:    burgundy.GetString("eos_host"),
		port:    burgundy.GetInt("eos_port"),
		chainID: burgundy.GetString("eos_chain"),
	}

	api, err := NewAPI(eosnet, burgundy.GetStringSlice("key"))
	if err != nil {
		return errors.Annotatef(err, "InitModule NewAPI failed")
	}
	crawler, err := NewCrawler(api)
	if err != nil {
		return errors.Annotatef(err, "InitModule NewCrawler failed")
	}
	crawlDuration := time.Millisecond * time.Duration(burgundy.GetInt("eos_crawl"))
	timer, err = NewTimer(crawler, crawlDuration, cancel)
	if err != nil {
		return errors.Annotatef(err, "InitModule NewTimer failed")
	}

	return nil
}
