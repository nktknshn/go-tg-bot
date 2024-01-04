package logging

import (
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
	"go.uber.org/zap"
)

type LoggerCreator = func(*zap.Logger) *zap.Logger

type ApplicationChatLoggerCreator = func(*zap.Logger, *telegram.TelegramUpdateContext) *zap.Logger
type ChatHandlerLoggerCreator = func(*zap.Logger, *telegram.TelegramUpdateContext) *zap.Logger

type TgbotLoggers struct {
	Base             *zap.Logger
	ChatsDistpatcher LoggerCreator
	ChatHandler      ChatHandlerLoggerCreator
	ApplicationChat  ApplicationChatLoggerCreator
	Component        LoggerCreator
}

func DefaultChatsDistpatcherLogger(logger *zap.Logger) *zap.Logger {
	return logger.Named("ChatsDistpatcher")
}

func DefaultChatHandlerLogger(logger *zap.Logger, tc *telegram.TelegramUpdateContext) *zap.Logger {
	return logger.Named("ChatHandler").With(zap.Int64("ChatID", tc.ChatID))
}

func DefaultComponentLogger(logger *zap.Logger) *zap.Logger {
	return logger.Named("Component")
}

func DefaultApplicationChat(logger *zap.Logger, tc *telegram.TelegramUpdateContext) *zap.Logger {
	return logger.Named("ApplicationChat").With(zap.Int64("ChatID", tc.ChatID))
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
