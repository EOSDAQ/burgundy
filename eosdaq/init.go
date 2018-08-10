package eosdaq

import (
	"burgundy/util"

	"go.uber.org/zap"
)

var mlog *zap.SugaredLogger

func init() {
	mlog, _ = util.InitLog("eosdaq", "console")
}
