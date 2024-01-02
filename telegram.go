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

// TelegramContext is a context related to a specific update
type TelegramContext struct {
	ChatID int64
	Ctx    context.Context
	Bot    TelegramBot
	Update BotUpdate
	Logger *zap.Logger
}

func (tc TelegramContext) AsTextMessage() (*TelegramContextTextMessage, bool) {
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
		TelegramContext: tc,
		Text:            m.Message,
		Message:         m,
	}, true
}

func (tc TelegramContext) AsCallback() (*TelegramContextCallback, bool) {
	u, ok := tc.Update.UpdateClass.(*tg.UpdateBotCallbackQuery)

	if !ok {
		return nil, false
	}

	return &TelegramContextCallback{
		TelegramContext:        tc,
		UpdateBotCallbackQuery: u,
	}, true
}

type TelegramUserChat struct {
	ChatID int64
	Chat   *tg.Chat
}

type TelegramContextTextMessage struct {
	TelegramContext
	Text    string
	Message *tg.Message
}

type TelegramContextCallback struct {
	TelegramContext
	UpdateBotCallbackQuery *tg.UpdateBotCallbackQuery
}

func (tc TelegramContext) AnswerCallbackQuery() {

	u, ok := tc.Update.UpdateClass.(*tg.UpdateBotCallbackQuery)

	if !ok {
		tc.Logger.Error("Update is not a callback query")
		return
	}

	tc.Bot.AnswerCallbackQuery(tc.Ctx, AnswerCallbackQueryParams{
		QueryID: u.QueryID,
	})
}
