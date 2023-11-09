package tgbot

import (
	"context"

	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

type ChatHandlerFactory func(*TelegramContext) ChatHandler

// ChatsDispatcher is a map of chats
// dispatches updates to chats
type ChatsDispatcher struct {
	Chats       map[int64]ChatHandler
	ChatFactory ChatHandlerFactory
}

type ChatsDispatcherProps struct {
	ChatFactory ChatHandlerFactory
}

func NewChatsDispatcher(props *ChatsDispatcherProps) *ChatsDispatcher {
	return &ChatsDispatcher{
		Chats:       make(map[int64]ChatHandler),
		ChatFactory: props.ChatFactory,
	}
}

func (cd *ChatsDispatcher) HandleUpdate(ctx context.Context, bot TelegramContextBot, update *models.Update) {

	var logger = GetLogger().With(
		zap.String("module", "ChatHandler"),
	)

	chatID := GetUpdateChatId(update)

	if chatID == 0 {
		logger.Debug("Update has no chat id, skipping", zap.Any("update", update))
		return
	}

	chat, ok := cd.Chats[chatID]

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
		logger.Debug("Creating new chat")
		chat = cd.ChatFactory(tc)
		cd.Chats[chatID] = chat
	}

	logger.Debug("Handling update")
	chat.HandleUpdate(tc)
}
