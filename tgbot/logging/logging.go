package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger() *zap.Logger {
	return DevLogger()
}

func FromCore(core zapcore.Core) *zap.Logger {
	return zap.New(core)
}

func DevLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()

	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.EncoderConfig.TimeKey = ""

	return zap.Must(cfg.Build())

}

func ProdLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()

	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.EncoderConfig.TimeKey = ""

	return zap.Must(cfg.Build())
}
