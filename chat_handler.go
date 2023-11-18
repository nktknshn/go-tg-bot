package tgbot

import (
	"context"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type CallbackAnswerer interface {
	AnswerCallbackQuery(context.Context, *bot.AnswerCallbackQueryParams) (bool, error)
}

type TelegramBot interface {
	CallbackAnswerer
	ChatRendererBot
}

type TelegramContext struct {
	ChatID int64
	Bot    TelegramBot
	Ctx    context.Context
	Update *models.Update
	Logger *zap.Logger
}

func (tc *TelegramContext) AnswerCallbackQuery() {
	tc.Bot.AnswerCallbackQuery(tc.Ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: tc.Update.CallbackQuery.ID,
		ShowAlert:       false,
	})
}

type ChatHandler interface {
	HandleUpdate(*TelegramContext)
}

type Handler[S any, A any, C any] struct {
	app        Application[S, A, C]
	appContext *ApplicationContext[S, A, C]
}

func NewHandler[S any, A any, C any](app Application[S, A, C], tc *TelegramContext) *Handler[S, A, C] {
	tc.Logger.Debug("NewHandler")

	tc.Logger.Debug("CreateAppState")
	appState := app.CreateAppState(tc)

	chatState := InternalChatState[S, A, C]{
		ChatID:           tc.ChatID,
		AppState:         appState,
		RenderedElements: []RenderedElement{},
		InputHandler:     nil,
		CallbackHandler:  nil,
		Renderer:         app.CreateChatRenderer(tc),
		TreeState:        nil,
		Lock:             &sync.Mutex{},
	}

	ac := &ApplicationContext[S, A, C]{
		App:    &app,
		State:  &chatState,
		Logger: GetLogger().With(zap.Int("ChatID", int(tc.ChatID))),
	}

	tc.Logger.Debug("PreRender")

	// prerender to get input handlers
	res := app.PreRender(ac)

	tc.Logger.Debug("New handler has been created.")

	return &Handler[S, A, C]{
		app: app,
		appContext: &ApplicationContext[S, A, C]{
			App:    &app,
			State:  &res.InternalChatState,
			Logger: ac.Logger,
		},
	}
}

func (h *Handler[S, A, C]) HandleUpdate(tc *TelegramContext) {
	tc.Logger.Debug("HandleUpdate")

	if tc.Update.Message != nil && tc.Update.Message.Text != "" {
		h.app.HandleMessage(h.appContext, tc)
		return
	}

	if tc.Update.CallbackQuery != nil {
		h.app.HandleCallback(h.appContext, tc)
		return
	}

	tc.Logger.Debug("Unkown update (neither message nor callback)")
}

func (a *Application[S, A, C]) NewHandler(tc *TelegramContext) *Handler[S, A, C] {
	return NewHandler[S, A, C](*a, tc)
}

func (a *Application[S, A, C]) ChatsDispatcher() *ChatsDispatcher {
	return NewChatsDispatcher(&ChatsDispatcherProps{
		ChatFactory: func(tc *TelegramContext) ChatHandler {
			return a.NewHandler(tc)
		},
	})
}
