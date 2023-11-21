package tgbot

import "go.uber.org/zap"

func GetLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()

	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.EncoderConfig.TimeKey = ""

	return zap.Must(cfg.Build())

}

var globalLogger = GetLogger()
