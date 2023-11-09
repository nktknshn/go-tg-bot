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
