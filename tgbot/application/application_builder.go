package application

import (
	"github.com/nktknshn/go-tg-bot/tgbot/render"
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
)

func (app *Application[S, C]) WithHandleActionExternal(f HandleActionFunc[S, C]) *Application[S, C] {
	app.HandleActionExternal = f
	return app
}

func (app *Application[S, C]) WithHandleAction(f HandleActionFunc[S, C]) *Application[S, C] {
	app.HandleAction = f
	return app
}

func (app *Application[S, C]) WithHandleMessage(f HandleMessageFunc[S, C]) *Application[S, C] {
	app.HandleMessage = f
	return app
}

func (app *Application[S, C]) WithHandleCallback(f HandleCallbackFunc[S, C]) *Application[S, C] {
	app.HandleCallback = f
	return app
}

func (app *Application[S, C]) WithHandleInit(f HandleInitFunc[S]) *Application[S, C] {
	app.HandleInit = f
	return app
}

func (app *Application[S, C]) WithRenderFunc(f RenderFuncType[S, C]) *Application[S, C] {
	app.RenderFunc = f
	return app
}

func (app *Application[S, C]) WithCreateRenderer(f func(*telegram.TelegramUpdateContext) render.ChatRenderer) *Application[S, C] {
	app.CreateChatRenderer = f
	return app
}

func (app *Application[S, C]) WithGlobalContext(f func(*ChatState[S, C]) C) *Application[S, C] {
	app.CreateGlobalContext = f
	return app
}

func (app *Application[S, C]) WithStateToComp(f StateToCompFuncType[S, C]) *Application[S, C] {
	app.StateToComp = f
	return app
}

func (app *Application[S, C]) WithCreateAppState(f func(*telegram.TelegramUpdateContext) S) *Application[S, C] {
	app.CreateAppState = f
	return app
}
