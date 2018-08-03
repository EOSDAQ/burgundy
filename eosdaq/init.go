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

type EosNet struct {
	host     string
	port     int
	contract string
}

func init() {
	mlog, _ = util.InitLog("eosdaq", "console")
}

func NewEosnet(host string, port int, contract string) *EosNet {
	return &EosNet{host, port, contract}
}

func InitModule(burgundy conf.ViperConfig, cancel <-chan os.Signal) error {

	eosnet := &EosNet{
		host:     burgundy.GetString("eos_host"),
		port:     burgundy.GetInt("eos_port"),
		contract: burgundy.GetString("eos_contract"),
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
