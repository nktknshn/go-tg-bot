package tgbot

import (
	"context"

	"github.com/nktknshn/go-tg-bot/tgbot/dispatcher"
	"github.com/nktknshn/go-tg-bot/tgbot/gotd"
	"github.com/nktknshn/go-tg-bot/tgbot/logging"
	"go.uber.org/zap"
)

type TgBot struct {
	chatsDispatcher *dispatcher.ChatsDispatcher
	logsSystem      logging.LogsSystem
}

func NewTgBot() *TgBot {
	return &TgBot{}
}

func (b *TgBot) logger() *zap.Logger {
	return b.logsSystem.Loggers().Tgbot()
}

func (b *TgBot) Run(ctx context.Context) error {
	b.logger().Debug("Starting telegram bot.")

	return gotd.Run(ctx, logging.DevLogger(), b.chatsDispatcher)
}
