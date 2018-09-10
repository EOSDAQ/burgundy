package models

import (
	"burgundy/conf"
	"burgundy/util"

	"go.uber.org/zap"
)

var mlog *zap.SugaredLogger

func init() {
	mlog, _ = util.InitLog("models", conf.Burgundy.GetString("loglevel"))
}
