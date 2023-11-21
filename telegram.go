package tgbot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"go.uber.org/zap"
)

type callbackAnswerer interface {
	AnswerCallbackQuery(context.Context, *bot.AnswerCallbackQueryParams) (bool, error)
}

// Interface for rendering messages into some interface (telegram, emulator, console, etc)
type TelegramBot interface {
	messageDeleter
	messageEditor
	messageSender
	callbackAnswerer
}

type TelegramContext struct {
	ChatID int64
	Bot    TelegramBot
	Ctx    context.Context
	Update *models.Update
	Logger *zap.Logger
}

func (tc TelegramContext) AnswerCallbackQuery() {
	tc.Bot.AnswerCallbackQuery(tc.Ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: tc.Update.CallbackQuery.ID,
		ShowAlert:       false,
	})
}
