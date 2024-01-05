package dispatcher

import (
	"context"
	"crypto/rand"
	"io"
	"sync"

	"github.com/nktknshn/go-tg-bot/helpers"
	"github.com/nktknshn/go-tg-bot/tgbot/logging"
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
	"go.uber.org/zap"
)

type ChatHandlerFactory interface {
	CreateChatHandler(*telegram.TelegramUpdateContext) ChatHandler
}

type ChatsDispatcherProps struct {
	ChatFactory ChatHandlerFactory
	loggers     logging.Loggers
}

// ChatsDispatcher is a map of chats
// dispatches updates to ChatHandlers
// TODO: AB testing
type ChatsDispatcher struct {
	chatHandlerFactory ChatHandlerFactory

	chatHandlers map[int64]ChatHandler

	logger  *zap.Logger
	loggers logging.Loggers

	chatLocks map[int64]*sync.Mutex

	stateLock *sync.Mutex

	randReader io.Reader
}

func NewChatsDispatcher(props *ChatsDispatcherProps) *ChatsDispatcher {

	loggers := props.loggers

	if props.loggers == nil {
		loggers = logging.NewLoggersDefault(logging.Logger())
	}

	logger := loggers.ChatsDispatcher()

	return &ChatsDispatcher{
		chatHandlers:       make(map[int64]ChatHandler),
		chatHandlerFactory: props.ChatFactory,
		logger:             logger,
		loggers:            loggers,
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

func (cd *ChatsDispatcher) createChatHandler(tc *telegram.TelegramUpdateContext) ChatHandler {
	chat := cd.chatHandlerFactory.CreateChatHandler(tc)
	cd.chatHandlers[tc.ChatID] = chat
	return chat
}

func (cd *ChatsDispatcher) createTelegramContext(ctx context.Context, bot telegram.TelegramBot, update telegram.BotUpdate) *telegram.TelegramUpdateContext {

	chatID := update.User.ID
	updateID, err := helpers.RandInt64(cd.randReader)
	logger := cd.loggers.Update(update, updateID)

	if err != nil {
		cd.logger.Error("Failed to generate update id. UpdateID will be 0", zap.Error(err))
	}

	return &telegram.TelegramUpdateContext{
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

// run in the goroutine
func (cd *ChatsDispatcher) handle(ctx context.Context, bot telegram.TelegramBot, update telegram.BotUpdate, chatLock *sync.Mutex) {

	tc := cd.createTelegramContext(ctx, bot, update)

	chatID := tc.ChatID
	logger := cd.logger.With(zap.Int64("ChatID", chatID))

	// ready for the next update
	defer chatLock.Unlock()

	chat, ok := cd.GetChatHandler(chatID)

	if !ok {
		logger.Debug("Creating new chat handler.")
		chat = cd.createChatHandler(tc)
	}

	chat.HandleUpdate(tc)

}

func (cd *ChatsDispatcher) HandleUpdate(ctx context.Context, bot telegram.TelegramBot, update telegram.BotUpdate) {

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
