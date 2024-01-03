package tgbot

import (
	"go.uber.org/zap"
)

// handle a certain chat update
type ChatHandler interface {
	HandleUpdate(*TelegramUpdateContext)
}

// updates must be put into the queue
type ChatHandlerImpl[S any, C any] struct {
	app     Application[S, C]
	appChat *ApplicationChat[S, C]
	logger  *zap.Logger
}

// Creates a new chat handler from an update
func NewChatHandler[S any, C any](app Application[S, C], tc *TelegramUpdateContext) *ChatHandlerImpl[S, C] {

	return &ChatHandlerImpl[S, C]{
		app:     app,
		appChat: NewApplicationChat[S, C](app, tc),
		logger: app.Loggers.
			ChatHandler(app.Loggers.Base).
			With(zap.Int64("ChatID", tc.ChatID)),
	}
}

func (h *ChatHandlerImpl[S, C]) HandleUpdate(tc *TelegramUpdateContext) {
	h.logger.Debug("HandleUpdate",
		zap.Any("UpdateType", tc.Update.UpdateClass.TypeName()),
		zap.Any("UpdateID", tc.UpdateID),
	)

	if tcm, ok := tc.AsTextMessage(); ok {
		h.app.HandleMessage(h.appChat, tcm)
		return
	}

	if tccb, ok := tc.AsCallback(); ok {
		h.app.HandleCallback(h.appChat, tccb)
		return
	}

	h.logger.Debug("Unkown update (neither message nor callback)")
}
