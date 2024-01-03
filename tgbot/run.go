package tgbot

import (
	"context"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

func Run(ctx context.Context, logger *zap.Logger, dispatcher *ChatsDispatcher) error {
	logger.Debug("Starting real telegram bot")

	handler := NewGotdHandler(dispatcher)

	opts := telegram.Options{
		Logger:        logger,
		UpdateHandler: handler,
	}

	err := telegram.BotFromEnvironment(ctx, opts,
		func(ctx context.Context, client *telegram.Client) error {

			api := tg.NewClient(client)

			sender := message.NewSender(api)
			handler.SetSender(sender)
			handler.SetClient(api)

			return nil
		},
		telegram.RunUntilCanceled,
	)

	return err

}
