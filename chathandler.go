package tgbot

import (
	"sync"

	"go.uber.org/zap"
)

type chatHandler interface {
	HandleUpdate(*TelegramContext)
}

type handler[S any, C any] struct {
	app        Application[S, C]
	appContext *ApplicationContext[S, C]
}

func NewHandler[S any, C any](app Application[S, C], tc *TelegramContext) *handler[S, C] {
	tc.Logger.Debug("NewHandler")

	tc.Logger.Debug("CreateAppState")
	appState := app.CreateAppState(tc)

	chatState := ChatState[S, C]{
		ChatID:           tc.ChatID,
		AppState:         appState,
		renderedElements: []RenderedElement{},
		inputHandler:     nil,
		callbackHandler:  nil,
		Renderer:         app.CreateChatRenderer(tc),
		treeState:        nil,
		lock:             &sync.Mutex{},
	}

	ac := &ApplicationContext[S, C]{
		App:    &app,
		State:  &chatState,
		Logger: GetLogger().With(zap.Int("ChatID", int(tc.ChatID))),
	}

	tc.Logger.Debug("PreRender")

	// prerender to get input handlers
	res := app.PreRender(ac)

	tc.Logger.Debug("New handler has been created.")

	return &handler[S, C]{
		app: app,
		appContext: &ApplicationContext[S, C]{
			App:    &app,
			State:  &res.InternalChatState,
			Logger: ac.Logger,
		},
	}
}

func (h *handler[S, C]) HandleUpdate(tc *TelegramContext) {
	tc.Logger.Debug("HandleUpdate")

	if tc, ok := tc.AsTextMessage(); ok {
		h.app.HandleMessage(h.appContext, tc)
		return
	} else if tc, ok := tc.AsCallback(); ok {
		h.app.HandleCallback(h.appContext, tc)
		return
	}

	tc.Logger.Debug("Unkown update (neither message nor callback)")
}

func (a *Application[S, C]) NewHandler(tc *TelegramContext) *handler[S, C] {
	return NewHandler[S, C](*a, tc)
}

func (a *Application[S, C]) ChatsDispatcher() *ChatsDispatcher {
	return NewChatsDispatcher(&ChatsDispatcherProps{
		ChatFactory: func(tc *TelegramContext) chatHandler {
			return a.NewHandler(tc)
		},
	})
}
