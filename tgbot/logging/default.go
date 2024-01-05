package logging

import (
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggersDefault struct {
	baseLogger *zap.Logger
}

type LoggersApplicationChatDefault struct {
	baseLogger *zap.Logger
}

func NewLoggersDefault(baseLogger *zap.Logger) *LoggersDefault {

	if baseLogger == nil {
		baseLogger = Logger()
	}

	return &LoggersDefault{
		baseLogger: baseLogger,
	}
}

func (ld *LoggersDefault) SetBase(base *zap.Logger) {
	ld.baseLogger = base
}

func (ld *LoggersDefault) Base() *zap.Logger {
	return ld.baseLogger
}

func (ld *LoggersDefault) SetFilter(filter FilterFunc) {
	ld.baseLogger = ld.baseLogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return NewFilteringCore(c, filter)
	}))
}

func (ld *LoggersDefault) IncreaseLevel(level zapcore.Level) {
	ld.baseLogger = ld.baseLogger.WithOptions(zap.IncreaseLevel(level))
}

func (ld *LoggersDefault) Update(update telegram.BotUpdate, updateID int64) *zap.Logger {
	return ld.baseLogger.Named(LoggerNameUpdate).With(
		zap.Int64("UpdateID", updateID),
		zap.Int64("ChatID", update.User.ID),
	)
}

func (ld *LoggersDefault) ChatsDispatcher() *zap.Logger {
	return ld.baseLogger.Named(LoggerNameChatsDispatcher)
}

func (ld *LoggersDefault) Tgbot() *zap.Logger {
	return ld.baseLogger.Named(LoggerNameTgbot)
}

func (ld *LoggersDefault) ApplicationChat(
	tuc *telegram.TelegramUpdateContext,
) LoggersApplicationChat {
	return &LoggersApplicationChatDefault{
		baseLogger: ld.baseLogger.Named(LoggerNameApplicationChat).With(
			zap.Int64("ChatID", tuc.ChatID),
		),
	}
}

func (ld *LoggersApplicationChatDefault) Init() *zap.Logger {
	return ld.baseLogger.Named(LoggerNameInit)
}

func (ld *LoggersApplicationChatDefault) Handler() *zap.Logger {
	return ld.baseLogger.Named(LoggerNameHandle)
}

func (ld *LoggersApplicationChatDefault) Action() *zap.Logger {
	return ld.baseLogger.Named(LoggerNameAction)
}

func (ld *LoggersApplicationChatDefault) Component() *zap.Logger {
	return ld.baseLogger.Named(LoggerNameComponent)
}

func (ld *LoggersApplicationChatDefault) Render() *zap.Logger {
	return ld.baseLogger.Named(LoggerNameRender)
}

func (ld *LoggersApplicationChatDefault) LockState() *zap.Logger {
	return ld.baseLogger.Named(LoggerNameLockState)
}
