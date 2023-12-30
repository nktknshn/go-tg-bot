package tgbot

import (
	"context"

	"go.uber.org/zap"
)

type ChatHandlerFactory func(*TelegramContext) chatHandler

type ChatsDispatcherProps struct {
	ChatFactory ChatHandlerFactory
}

// ChatsDispatcher is a map of chats
// dispatches updates to chats
type ChatsDispatcher struct {
	ChatHandlers       map[int64]chatHandler
	ChatHandlerFactory ChatHandlerFactory
	Logger             *zap.Logger
}

func NewChatsDispatcher(props *ChatsDispatcherProps) *ChatsDispatcher {
	return &ChatsDispatcher{
		ChatHandlers:       make(map[int64]chatHandler),
		ChatHandlerFactory: props.ChatFactory,
		Logger:             GetLogger(),
	}
}

func (cd *ChatsDispatcher) newTelegramContextLogger(bot TelegramBot, chatID int64, update BotUpdate) *zap.Logger {

	return GetLogger().With(
		zap.Int64("chatID", chatID),
		// zap.Int64("updateID", update.ID),
	)
}

func (cd *ChatsDispatcher) HandleUpdate(ctx context.Context, bot TelegramBot, update BotUpdate) error {

	cd.Logger.Debug("HandleUpdate", zap.Any("update", update))

	// logger := cd.Logger.With(zap.Int64("updateID", update.ID))
	logger := cd.Logger

	chatID := update.User.ID

	if chatID == 0 {
		logger.Debug("Update has no chat id, skipping.")
		return nil
	}

	chat, ok := cd.ChatHandlers[chatID]

	tc := &TelegramContext{
		Bot:    bot,
		Logger: cd.newTelegramContextLogger(bot, chatID, update),
		Ctx:    ctx,
		ChatID: chatID,
		Update: update,
	}

	if !ok {
		logger.Debug("Creating new chat handler")
		chat = cd.ChatHandlerFactory(tc)
		cd.ChatHandlers[chatID] = chat
	}

	logger.Debug("Handling update")
	chat.HandleUpdate(tc)

	return nil
}
