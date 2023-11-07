package tgbot

import "go.uber.org/zap"

func GetLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()

	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	return zap.Must(cfg.Build())

}
