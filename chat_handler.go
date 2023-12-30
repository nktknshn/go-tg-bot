package tgbot

import (
	"sync"

	"go.uber.org/zap"
)

type ChatHandler interface {
	HandleUpdate(*TelegramContext)
}

type ChatHandlerImpl[S any, C any] struct {
	app        Application[S, C]
	appContext *ApplicationContext[S, C]
}

func NewHandler[S any, C any](app Application[S, C], tc *TelegramContext) *ChatHandlerImpl[S, C] {
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

	return &ChatHandlerImpl[S, C]{
		app: app,
		appContext: &ApplicationContext[S, C]{
			App:    &app,
			State:  &res.InternalChatState,
			Logger: ac.Logger,
		},
	}
}

func (h *ChatHandlerImpl[S, C]) HandleUpdate(tc *TelegramContext) {
	tc.Logger.Debug("HandleUpdate", zap.Any("update", tc))

	if tcm, ok := tc.AsTextMessage(); ok {
		h.app.HandleMessage(h.appContext, tcm)
		return
	}

	if tccb, ok := tc.AsCallback(); ok {
		h.app.HandleCallback(
			h.appContext,
			tccb,
		)
		return
	}

	tc.Logger.Debug("Unkown update (neither message nor callback)")
}

func (a *Application[S, C]) NewHandler(tc *TelegramContext) *ChatHandlerImpl[S, C] {
	return NewHandler[S, C](*a, tc)
}

func (a *Application[S, C]) ChatsDispatcher() *ChatsDispatcher {
	return NewChatsDispatcher(&ChatsDispatcherProps{
		ChatFactory: func(tc *TelegramContext) ChatHandler {
			return a.NewHandler(tc)
		},
	})
}
