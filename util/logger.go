package util

import (
	"strings"
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
func InitLog(name string, env string) (log *zap.SugaredLogger, err error) {

	cfg := zap.NewDevelopmentConfig()
	cfg.Encoding = "console"
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	enccfg := zap.NewDevelopmentEncoderConfig()

	switch strings.ToLower(env) {
	case "info":
		cfg.Encoding = "json"
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		enccfg = zap.NewDevelopmentEncoderConfig()
	case "error":
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "json"
		cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		enccfg = zap.NewProductionEncoderConfig()
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
