package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func InitLogger() {
	var cfg zap.Config

	outputLevel := zapcore.InfoLevel

	cfg = zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}
	cfg.InitialFields = map[string]interface{}{"name": "kabutar"}
	cfg.Encoding = os.Getenv("LOGGING")
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Level = zap.NewAtomicLevelAt(outputLevel)

	simple, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	logger = simple.Sugar()
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Errof(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}
