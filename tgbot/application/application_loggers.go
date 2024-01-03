package application

import "github.com/nktknshn/go-tg-bot/tgbot/logging"

func (app *Application[S, C]) SetLoggers(loggers logging.TgbotLoggers) {
	app.Loggers = loggers
}
