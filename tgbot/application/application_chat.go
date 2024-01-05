package application

import (
	"sync"

	"github.com/nktknshn/go-tg-bot/tgbot/logging"
	"github.com/nktknshn/go-tg-bot/tgbot/render"
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
)

// tie together an application methods and a chat state
type ApplicationChat[S any, C any] struct {
	App   *Application[S, C]
	State *ApplicationChatState[S, C]

	// loggers for different parts of the app
	Loggers logging.LoggersApplicationChat
}

func (ac *ApplicationChat[S, C]) SetChatState(chatState *ApplicationChatState[S, C]) {
	ac.State = chatState
}

func NewApplicationChat[S any, C any](app *Application[S, C], tc *telegram.TelegramUpdateContext) *ApplicationChat[S, C] {

	chatLogger := app.Loggers.ApplicationChat(tc)

	appState := app.CreateAppState(app, tc, chatLogger.Init())

	chatState := ApplicationChatState[S, C]{
		ChatID:           tc.ChatID,
		AppState:         appState,
		renderedElements: []render.RenderedElement{},
		inputHandler:     nil,
		callbackHandler:  nil,
		treeState:        nil,
		Renderer:         app.CreateChatRenderer(tc),
		lock:             &sync.Mutex{},
	}

	// compute the first state
	res := app.ComputeNextState(&chatState, ComputeNextStateProps{
		Logger: chatLogger.Init().Named(logging.LoggerNameComponent),
	})

	return &ApplicationChat[S, C]{
		App:     app,
		State:   &res.NextChatState,
		Loggers: chatLogger,
	}
}
