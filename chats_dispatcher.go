package tgbot

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

type ChatHandlerFactory interface {
	CreateChatHandler(*TelegramContext) ChatHandler
}

type ChatsDispatcherProps struct {
	ChatFactory ChatHandlerFactory
}

// ChatsDispatcher is a map of chats
// dispatches updates to ChatHandlers
// TODO: AB testing

type ChatsDispatcher struct {
	chatHandlerFactory ChatHandlerFactory

	chatHandlers map[int64]ChatHandler
	logger       *zap.Logger

	chatLocks map[int64]*sync.Mutex

	stateLock *sync.Mutex
}

func NewChatsDispatcher(props *ChatsDispatcherProps) *ChatsDispatcher {
	return &ChatsDispatcher{
		chatHandlers:       make(map[int64]ChatHandler),
		chatHandlerFactory: props.ChatFactory,
		logger:             GetLogger(),
		chatLocks:          make(map[int64]*sync.Mutex),
		stateLock:          &sync.Mutex{},
	}
}

func (cd *ChatsDispatcher) SetLogger(logger *zap.Logger) {
	cd.logger = logger
}

func (cd *ChatsDispatcher) newTelegramContextLogger(bot TelegramBot, chatID int64, update BotUpdate) *zap.Logger {

	return cd.logger.With(
		zap.Int64("chatID", chatID),
		// zap.Int64("updateID", update.ID),
	)
}

func (cd *ChatsDispatcher) createChatHandler(tc *TelegramContext) ChatHandler {
	chatID := tc.ChatID
	chat := cd.chatHandlerFactory.CreateChatHandler(tc)

	cd.chatHandlers[chatID] = chat

	return chat
}

func (cd *ChatsDispatcher) createTelegramContext(ctx context.Context, bot TelegramBot, update BotUpdate) *TelegramContext {

	chatID := update.User.ID
	logger := cd.newTelegramContextLogger(bot, chatID, update)

	return &TelegramContext{
		Bot:    bot,
		Logger: logger,
		ChatID: chatID,
		Update: update,
	}
}

func (cd *ChatsDispatcher) getChatLock(chatID int64) *sync.Mutex {

	if _, ok := cd.chatLocks[chatID]; !ok {
		cd.chatLocks[chatID] = &sync.Mutex{}
	}

	return cd.chatLocks[chatID]
}

func (cd *ChatsDispatcher) GetChatHandler(chatID int64) (ChatHandler, bool) {
	chat, ok := cd.chatHandlers[chatID]
	return chat, ok
}

func (cd *ChatsDispatcher) handle(ctx context.Context, bot TelegramBot, update BotUpdate, chatLock *sync.Mutex) {

	tc := cd.createTelegramContext(ctx, bot, update)

	chatID := tc.ChatID
	logger := tc.Logger

	defer chatLock.Unlock()

	chat, ok := cd.GetChatHandler(chatID)

	if !ok {
		logger.Debug("Creating new chat handler")
		chat = cd.createChatHandler(tc)
	}

	cd.logger.Debug("Handling update")

	chat.HandleUpdate(tc)

}

func (cd *ChatsDispatcher) HandleUpdate(ctx context.Context, bot TelegramBot, update BotUpdate) {

	logger := cd.logger
	chatID := update.User.ID

	if chatID == 0 {
		logger.Debug("Update has no chat id, skipping.")
		return
	}

	cd.stateLock.Lock()

	chatLock := cd.getChatLock(chatID)
	chatLock.Lock()

	go cd.handle(ctx, bot, update, chatLock)

	cd.stateLock.Unlock()

}
