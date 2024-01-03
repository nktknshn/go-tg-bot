package telegram

import (
	"context"

	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

// TelegramUpdateContext is a context related to a specific update
type TelegramUpdateContext struct {
	Ctx context.Context

	ChatID int64

	// for tracking purposes
	UpdateID int64
	Update   BotUpdate

	Bot TelegramBot

	UpdateLogger *zap.Logger
}

func (tc TelegramUpdateContext) AsTextMessage() (*TelegramContextTextMessage, bool) {
	u, ok := tc.Update.UpdateClass.(*tg.UpdateNewMessage)

	if !ok {
		return nil, false
	}

	if u.Message == nil {
		return nil, false
	}

	m, ok := u.Message.(*tg.Message)

	if !ok {
		return nil, false
	}

	return &TelegramContextTextMessage{
		TelegramUpdateContext: tc,
		Text:                  m.Message,
		Message:               m,
	}, true
}

type TelegramContextTextMessage struct {
	TelegramUpdateContext
	Text    string
	Message *tg.Message
}

type TelegramContextCallback struct {
	TelegramUpdateContext
	UpdateBotCallbackQuery *tg.UpdateBotCallbackQuery
}

func (tc TelegramContextCallback) AnswerCallbackQuery() {

	u, ok := tc.Update.UpdateClass.(*tg.UpdateBotCallbackQuery)

	if !ok {
		tc.UpdateLogger.Error("Update is not a callback query")
		return
	}

	tc.Bot.AnswerCallbackQuery(tc.Ctx, AnswerCallbackQueryParams{
		QueryID: u.QueryID,
	})
}

func (tc TelegramUpdateContext) AsCallback() (*TelegramContextCallback, bool) {
	u, ok := tc.Update.UpdateClass.(*tg.UpdateBotCallbackQuery)

	if !ok {
		return nil, false
	}

	return &TelegramContextCallback{
		TelegramUpdateContext:  tc,
		UpdateBotCallbackQuery: u,
	}, true
}
