package tgbot

import "go.uber.org/zap"

type LoggerCreator = func(*zap.Logger) *zap.Logger

type TgbotLoggers struct {
	Base             *zap.Logger
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

func DefaultApplicationChat(logger *zap.Logger) *zap.Logger {
	return logger.Named("ApplicationChat")
}

var DefaultLoggers = TgbotLoggers{
	Base:             DevLogger(),
	ChatsDistpatcher: DefaultChatsDistpatcherLogger,
	ChatHandler:      DefaultChatHandlerLogger,
	Component:        DefaultComponentLogger,
	ApplicationChat:  DefaultApplicationChat,
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

var globalLogger = DevLogger()
