package tgbot

import "go.uber.org/zap"

type LoggerCreator = func(*zap.Logger) *zap.Logger

type TgbotLoggers struct {
	ChatsDistpatcher LoggerCreator
	ChatHandler      LoggerCreator
	Component        LoggerCreator
	ApplicationChat  LoggerCreator
}

func DefaultChatsDistpatcherLogger(logger *zap.Logger) *zap.Logger {
	return logger.Named("ChatsDistpatcher")
}

func DefaultChatHandlerLogger(logger *zap.Logger) *zap.Logger {
	return logger.Named("ChatHandler")
}

func DefaultComponentLogger(logger *zap.Logger) *zap.Logger {
	return logger.Named("Component")
}

var DefaultLoggers = TgbotLoggers{
	ChatsDistpatcher: DefaultChatsDistpatcherLogger,
	ChatHandler:      DefaultChatHandlerLogger,
	Component:        DefaultComponentLogger,
}

func GetLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()

	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.EncoderConfig.TimeKey = ""

	return zap.Must(cfg.Build())

}

var globalLogger = GetLogger()
