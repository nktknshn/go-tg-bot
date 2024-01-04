package logging

import "go.uber.org/zap"

type LoggerCreator = func(*zap.Logger) *zap.Logger

type TgbotLoggers struct {
	Base             *zap.Logger
	ChatsDistpatcher LoggerCreator
	ChatHandler      LoggerCreator
	ApplicationChat  func(*zap.Logger, int64) *zap.Logger
	Component        LoggerCreator
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

func DefaultApplicationChat(logger *zap.Logger, chatID int64) *zap.Logger {
	return logger.Named("ApplicationChat").With(zap.Int64("ChatID", chatID))
}

var DefaultLoggers = TgbotLoggers{
	Base:             DevLogger(),
	ChatsDistpatcher: DefaultChatsDistpatcherLogger,
	ChatHandler:      DefaultChatHandlerLogger,
	Component:        DefaultComponentLogger,
	ApplicationChat:  DefaultApplicationChat,
}

func Logger() *zap.Logger {
	return DevLogger()
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
