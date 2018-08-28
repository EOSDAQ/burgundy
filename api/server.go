package main

import (
	conf "burgundy/conf"
	"burgundy/crawler"
	"burgundy/util"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
)

func prepareServer(burgundy *conf.ViperConfig, cancel <-chan os.Signal, db *gorm.DB) bool {

	log, err := util.InitLog("server", burgundy.GetString("logmode"))
	if err != nil {
		log.Infow("InitLog", "err", errors.Details(err))
		return false
	}

	if err := crawler.InitModule(burgundy, cancel, db); err != nil {
		log.Infow("InitModule eosdaq", "err", errors.Details(err))
		return false
	}

	return true
}
