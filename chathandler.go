package tgbot

import (
	"context"

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

type ChatHandler interface {
	HandleUpdate(*TelegramContext)
}

type Handler[S any, A any] struct {
	app        Application[S, A]
	appContext *ApplicationContext[S, A]
}

func NewHandler[S any, A any](app Application[S, A], tc *TelegramContext) *Handler[S, A] {
	tc.Logger.Debug("NewHandler")

	tc.Logger.Debug("CreateAppState")
	appState := app.CreateAppState(tc)

	chatState := InternalChatState[S, A]{
		ChatID:           tc.ChatID,
		AppState:         appState,
		RenderedElements: []RenderedElement{},
		InputHandler:     nil,
		CallbackHandler:  nil,
		Renderer:         app.CreateChatRenderer(tc),
		TreeState:        nil,
	}

	ac := &ApplicationContext[S, A]{
		App:    &app,
		State:  &chatState,
		Logger: GetLogger().With(zap.Int("ChatID", int(tc.ChatID))),
	}

	tc.Logger.Debug("PreRender")
	res := app.PreRender(ac)

	tc.Logger.Debug("New handler has been created.")

	return &Handler[S, A]{
		app: app,
		appContext: &ApplicationContext[S, A]{
			App:    &app,
			State:  &res.InternalChatState,
			Logger: ac.Logger,
		},
	}
}

func (h *Handler[S, A]) HandleUpdate(tc *TelegramContext) {
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
