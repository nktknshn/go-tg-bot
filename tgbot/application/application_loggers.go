package application

import (
	"github.com/nktknshn/go-tg-bot/tgbot/logging"
	"go.uber.org/zap"
)

type ApplicationChatLoggers struct {
	Root      *zap.Logger
	Init      *zap.Logger
	Component *zap.Logger
	Handle    *zap.Logger
	Action    *zap.Logger
	Render    *zap.Logger
}

func (app *Application[S, C]) SetLoggers(loggers logging.Loggers) {
	app.Loggers = loggers
}
