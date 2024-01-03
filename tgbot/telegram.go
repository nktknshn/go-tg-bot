package tgbot

import (
	"context"

	"github.com/gotd/td/tg"

	"go.uber.org/zap"
)

type AnswerCallbackQueryParams struct {
	QueryID int64
}

type CallbackAnswerer interface {
	AnswerCallbackQuery(context.Context, AnswerCallbackQueryParams) (bool, error)
}

// Interface for rendering messages into some interface (telegram, emulator, console, etc)
type TelegramBot interface {
	MessageDeleter
	MessageEditor
	MessageSender
	CallbackAnswerer
}

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

type TelegramUserChat struct {
	ChatID int64
	Chat   *tg.Chat
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
