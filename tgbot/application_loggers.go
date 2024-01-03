package tgbot

func (app *Application[S, C]) SetLoggers(loggers TgbotLoggers) {
	app.Loggers = loggers
}
