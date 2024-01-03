package tgbot

import (
	"github.com/nktknshn/go-tg-bot/tgbot/application"
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
	"go.uber.org/zap"
)

// handle a certain chat update
type ChatHandler interface {
	HandleUpdate(*telegram.TelegramUpdateContext)
}

// updates must be put into the queue
type ChatHandlerImpl[S any, C any] struct {
	app     application.Application[S, C]
	appChat *application.ApplicationChat[S, C]
	logger  *zap.Logger
}

// Creates a new chat handler from an update
func NewChatHandler[S any, C any](app application.Application[S, C], tc *telegram.TelegramUpdateContext) *ChatHandlerImpl[S, C] {

	return &ChatHandlerImpl[S, C]{
		app:     app,
		appChat: application.NewApplicationChat[S, C](app, tc),
		logger: app.Loggers.
			ChatHandler(app.Loggers.Base).
			With(zap.Int64("ChatID", tc.ChatID)),
	}
}

func (h *ChatHandlerImpl[S, C]) HandleUpdate(tc *telegram.TelegramUpdateContext) {
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

// func (a *Application[S, C]) NewHandler(tc *telegram.TelegramUpdateContext) *ChatHandlerImpl[S, C] {
// 	return NewChatHandler[S, C](*a, tc)
// }

// func (a *Application[S, C]) ChatsDispatcher() *ChatsDispatcher {

// 	return NewChatsDispatcher(&ChatsDispatcherProps{
// 		ChatFactory: &factoryFunc{
// 			f: func(tc *telegram.TelegramUpdateContext) ChatHandler {
// 				return a.NewHandler(tc)
// 			},
// 		},
// 	})
// }

// type factoryFunc struct {
// 	f func(*telegram.TelegramUpdateContext) ChatHandler
// }

// func (f *factoryFunc) CreateChatHandler(tc *telegram.TelegramUpdateContext) ChatHandler {
// 	return f.f(tc)
// }
