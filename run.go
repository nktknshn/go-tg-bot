package tgbot

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func RunReal(logger *zap.Logger, dispatcher *ChatsDispatcher) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	token := os.Getenv("BOT_TOKEN")

	if token == "" {
		logger.Fatal("BOT_TOKEN env variable is not set")
		os.Exit(1)
	}

	bot, err := bot.New(token, bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		dispatcher.HandleUpdate(ctx, bot, update)
	}))

	if err != nil {
		logger.Fatal("Error creating bot", zap.Error(err))
		os.Exit(1)
	}

	bot.Start(ctx)

}
