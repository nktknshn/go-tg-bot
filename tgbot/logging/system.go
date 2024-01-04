package logging

import (
	"fmt"
	"os"
	"path"

	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogsSystem struct {
	Base       *zap.Logger
	LogsFolder string

	logger *zap.Logger
}

func NewLogsSystem(base *zap.Logger, logsFolder string) *LogsSystem {
	return &LogsSystem{
		Base:       base,
		LogsFolder: logsFolder,
		logger:     base.Named(LoggerNameLogsSystem),
	}
}

func (ul *LogsSystem) ApplicationChatLogger(
	l *zap.Logger,
	tc *telegram.TelegramUpdateContext,
) *zap.Logger {

	fname := fmt.Sprintf("user_%d.log", tc.ChatID)
	chatLogFile := path.Join(ul.LogsFolder, fname)

	ul.logger.Info("Opening log file", zap.String("path", chatLogFile))

	file, err := os.OpenFile(chatLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	l = l.Named(LoggerNameApplicationChat)
	l = l.WithOptions(
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {

			if err != nil {
				ul.logger.Error("failed to open log file", zap.Error(err))
				return core
			}

			return zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(zapcore.Lock(file)),
				zap.DebugLevel,
			)
		}),
	)

	return l
}

func (ul *LogsSystem) TgbotLoggers() *TgbotLoggers {

	return &TgbotLoggers{
		Base: ul.Base,
		ChatsDistpatcher: func(l *zap.Logger) *zap.Logger {
			return l.Named("ChatsDistpatcher")
		},
		ChatHandler: func(l *zap.Logger, tc *telegram.TelegramUpdateContext) *zap.Logger {
			return DefaultChatHandlerLogger(l, tc)
		},
		ApplicationChat: func(l *zap.Logger, tc *telegram.TelegramUpdateContext) *zap.Logger {
			return ul.ApplicationChatLogger(l, tc)
		},
		Component: func(l *zap.Logger) *zap.Logger {
			return l.Named(LoggerNameComponent).WithOptions()
		},
	}
}
