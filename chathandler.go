package tgbot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type TelegramContext struct {
	ChatID int64
	Bot    *bot.Bot
	Ctx    context.Context
	Update *models.Update
	Logger *zap.Logger
}

type ChatHandler interface {
	HandleUpdate(*TelegramContext)
}
