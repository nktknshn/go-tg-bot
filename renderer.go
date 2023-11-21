package tgbot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var ErrMessageNotFound = fmt.Errorf("message not found")

type messageDeleter interface {
	DeleteMessage(ctx context.Context, params *bot.DeleteMessageParams) (bool, error)
}

type messageEditor interface {
	EditMessageText(ctx context.Context, params *bot.EditMessageTextParams) (*models.Message, error)
}

type messageSender interface {
	SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error)
}
