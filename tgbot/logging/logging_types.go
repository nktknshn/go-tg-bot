package logging

import (
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
	"go.uber.org/zap"
)

type Loggers interface {
	SetBase(*zap.Logger)
	SetFilter(FilterFunc)
	Base() *zap.Logger
	Update(telegram.BotUpdate, int64) *zap.Logger
	ChatsDispatcher() *zap.Logger
	Tgbot() *zap.Logger
	ApplicationChat(
		*telegram.TelegramUpdateContext,
	) LoggersApplicationChat
}

type LoggersApplicationChat interface {
	Init() *zap.Logger
	Handler() *zap.Logger
	Action() *zap.Logger
	Component() *zap.Logger
	Render() *zap.Logger
}
