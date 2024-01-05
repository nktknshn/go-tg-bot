package dispatcher

import (
	"github.com/nktknshn/go-tg-bot/tgbot/application"
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
)

func ForApplication[S, C any](app *application.Application[S, C]) *ChatsDispatcher {
	return NewChatsDispatcher(&ChatsDispatcherProps{
		ChatFactory: &factoryFromFunc{
			f: func(tc *telegram.TelegramUpdateContext) ChatHandler {
				return application.NewApplicationChat[S, C](
					app,
					tc,
				)
			},
		},
	})
}

type factoryFromFunc struct {
	f func(*telegram.TelegramUpdateContext) ChatHandler
}

func (f *factoryFromFunc) CreateChatHandler(tc *telegram.TelegramUpdateContext) ChatHandler {
	return f.f(tc)
}
