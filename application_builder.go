package tgbot

func (app *Application[S, C]) WithHandleActionExternal(f handleActionFunc[S, C]) *Application[S, C] {
	app.HandleActionExternal = f
	return app
}

func (app *Application[S, C]) WithHandleAction(f handleActionFunc[S, C]) *Application[S, C] {
	app.HandleAction = f
	return app
}

func (app *Application[S, C]) WithHandleMessage(f handleMessageFunc[S, C]) *Application[S, C] {
	app.HandleMessage = f
	return app
}

func (app *Application[S, C]) WithHandleCallback(f handleCallbackFunc[S, C]) *Application[S, C] {
	app.HandleCallback = f
	return app
}

func (app *Application[S, C]) WithHandleInit(f handleInitFunc[S]) *Application[S, C] {
	app.HandleInit = f
	return app
}

func (app *Application[S, C]) WithRenderFunc(f renderFuncType[S, C]) *Application[S, C] {
	app.RenderFunc = f
	return app
}

func (app *Application[S, C]) WithCreateRenderer(f func(*TelegramContext) ChatRenderer) *Application[S, C] {
	app.CreateChatRenderer = f
	return app
}

func (app *Application[S, C]) WithCreateGlobalContext(f func(*ChatState[S, C]) C) *Application[S, C] {
	app.CreateGlobalContext = f
	return app
}

func (app *Application[S, C]) WithStateToComp(f stateToCompFuncType[S, C]) *Application[S, C] {
	app.StateToComp = f
	return app
}

func (app *Application[S, C]) WithCreateAppState(f func(*TelegramContext) S) *Application[S, C] {
	app.CreateAppState = f
	return app
}
