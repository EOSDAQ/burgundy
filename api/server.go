package main

import (
	conf "burgundy/conf"
	"burgundy/eosdaq"
	"burgundy/util"
	"os"

	"github.com/juju/errors"
)

func prepareServer(burgundy conf.ViperConfig, cancel <-chan os.Signal) bool {

	log, err := util.InitLog("server", burgundy.GetString("logmode"))
	if err != nil {
		log.Infow("InitLog", "err", errors.Details(err))
		return false
	}

	if err := eosdaq.InitModule(burgundy, cancel); err != nil {
		log.Infow("InitModule eosdaq", "err", errors.Details(err))
		return false
	}

	return true
}
