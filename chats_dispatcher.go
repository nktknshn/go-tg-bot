package tgbot

import (
	"context"
	"crypto/rand"
	"io"
	"sync"

	"go.uber.org/zap"
)

type ChatHandlerFactory interface {
	CreateChatHandler(*TelegramUpdateContext) ChatHandler
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

	randReader io.Reader
}

func NewChatsDispatcher(props *ChatsDispatcherProps) *ChatsDispatcher {
	return &ChatsDispatcher{
		chatHandlers:       make(map[int64]ChatHandler),
		chatHandlerFactory: props.ChatFactory,
		logger:             DevLogger(),
		chatLocks:          make(map[int64]*sync.Mutex),
		stateLock:          &sync.Mutex{},
		randReader:         rand.Reader,
	}
}

func (cd *ChatsDispatcher) ResetChats() {
	cd.stateLock.Lock()
	cd.chatHandlers = make(map[int64]ChatHandler)
	cd.chatLocks = make(map[int64]*sync.Mutex)
	cd.stateLock.Unlock()
}

func (cd *ChatsDispatcher) ResetChat(chatID int64) {
	cd.stateLock.Lock()
	delete(cd.chatHandlers, chatID)
	delete(cd.chatLocks, chatID)
	cd.stateLock.Unlock()
}

func (cd *ChatsDispatcher) SetLogger(logger *zap.Logger) {
	cd.logger = logger
}

func (cd *ChatsDispatcher) newTelegramContextLogger(bot TelegramBot, chatID int64, update BotUpdate, updateID int64) *zap.Logger {

	return DevLogger().
		Named("TelegramUpdateContext").
		With(
			zap.Int64("chatID", chatID),
			zap.Int64("updateID", updateID),
		)
}

func (cd *ChatsDispatcher) createChatHandler(tc *TelegramUpdateContext) ChatHandler {
	chat := cd.chatHandlerFactory.CreateChatHandler(tc)
	cd.chatHandlers[tc.ChatID] = chat
	return chat
}

func (cd *ChatsDispatcher) createTelegramContext(ctx context.Context, bot TelegramBot, update BotUpdate) *TelegramUpdateContext {

	chatID := update.User.ID
	updateID, err := RandInt64(cd.randReader)
	logger := cd.newTelegramContextLogger(bot, chatID, update, updateID)

	if err != nil {
		cd.logger.Error("Failed to generate update id. UpdateID will be 0", zap.Error(err))
	}

	return &TelegramUpdateContext{
		Ctx:          ctx,
		Bot:          bot,
		UpdateLogger: logger,
		ChatID:       chatID,
		Update:       update,
		UpdateID:     updateID,
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
	logger := cd.logger.With(zap.Int64("ChatID", chatID))

	defer chatLock.Unlock()

	chat, ok := cd.GetChatHandler(chatID)

	if !ok {
		logger.Debug("Creating new chat handler.")
		chat = cd.createChatHandler(tc)
	}

	chat.HandleUpdate(tc)

}

func (cd *ChatsDispatcher) HandleUpdate(ctx context.Context, bot TelegramBot, update BotUpdate) {

	logger := cd.logger
	chatID := update.User.ID

	logger.Debug("Handling update", zap.Int64("ChatID", chatID), zap.Any("UpdateType", update.UpdateClass.TypeName()))

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
