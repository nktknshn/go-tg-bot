package tgbot

import (
	"go.uber.org/zap"
)

// handle a certain chat update
type ChatHandler interface {
	HandleUpdate(*TelegramContext)
}

// updates must be put into the queue
type ChatHandlerImpl[S any, C any] struct {
	app     Application[S, C]
	appChat *ApplicationChat[S, C]
}

// Creates a new chat handler from an update
func NewChatHandler[S any, C any](app Application[S, C], tc *TelegramContext) *ChatHandlerImpl[S, C] {
	tc.Logger.Debug("New handler has been created.")

	return &ChatHandlerImpl[S, C]{
		app:     app,
		appChat: NewApplicationChat[S, C](app, tc),
	}
}

func (h *ChatHandlerImpl[S, C]) HandleUpdate(tc *TelegramContext) {
	tc.Logger.Debug("HandleUpdate", zap.Any("update", tc))

	if tcm, ok := tc.AsTextMessage(); ok {
		h.app.HandleMessage(h.appChat, tcm)
		return
	}

	if tccb, ok := tc.AsCallback(); ok {
		h.app.HandleCallback(
			h.appChat,
			tccb,
		)
		return
	}

	tc.Logger.Debug("Unkown update (neither message nor callback)")
}
