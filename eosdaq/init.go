package eosdaq

import (
	"burgundy/conf"
	"burgundy/util"

	"go.uber.org/zap"
)

var mlog *zap.SugaredLogger

func init() {
	mlog, _ = util.InitLog("eosdaq", conf.Burgundy.GetString("loglevel"))
}
