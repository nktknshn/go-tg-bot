package tgbot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type ChatHandlerFactory func(*TelegramContext) ChatHandler

// ChatsDispatcher is a map of chats
// dispatches updates to chats
type ChatsDispatcher struct {
	chats       map[int64]ChatHandler
	chatFactory ChatHandlerFactory
}

func NewChatsDispatcher(chatFactory ChatHandlerFactory) *ChatsDispatcher {
	return &ChatsDispatcher{
		chats:       make(map[int64]ChatHandler),
		chatFactory: chatFactory,
	}
}

func (cd *ChatsDispatcher) HandleUpdate(ctx context.Context, bot *bot.Bot, update *models.Update) {

	var logger = GetLogger().With(
		zap.String("module", "ChatHandler"),
	)

	chatID := GetUpdateChatId(update)

	if chatID == 0 {
		logger.Debug("Update has no chat id, skipping", zap.Any("update", update))
		return
	}

	chat, ok := cd.chats[chatID]
	tc := &TelegramContext{
		Bot: bot,
		Logger: logger.With(
			zap.Int64("chat_id", chatID),
			zap.Any("update", update),
		),
		Ctx:    ctx,
		ChatID: chatID,
		Update: update,
	}

	if !ok {
		chat = cd.chatFactory(tc)
		cd.chats[chatID] = chat
	}

	chat.HandleUpdate(tc)
}
