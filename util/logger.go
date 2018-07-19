package util

import (
	"time"

	"github.com/juju/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TimeEncoder for logging time format.
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.00"))
}

var mlog *zap.SugaredLogger

// InitLog returns logger instance.
func InitLog(name string, logMode string) (log *zap.SugaredLogger, err error) {

	var cfg zap.Config
	var enccfg zapcore.EncoderConfig
	if logMode == "json" {
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "json"
		enccfg = zap.NewProductionEncoderConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.Encoding = "console"
		enccfg = zap.NewDevelopmentEncoderConfig()
	}
	enccfg.EncodeTime = TimeEncoder
	enccfg.CallerKey = ""
	enccfg.LevelKey = ""
	cfg.EncoderConfig = enccfg

	logger, err := cfg.Build()
	if err != nil {
		return nil, errors.Annotatef(err, "InitLog")
	}
	defer logger.Sync()

	mlog = logger.Sugar().Named(name)

	return mlog, nil
}
