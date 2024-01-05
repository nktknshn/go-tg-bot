package control

import "github.com/nktknshn/go-tg-bot/tgbot/dispatcher"

type TgBotControl interface {
	ResetChats()
	ResetChat(chatID int64)
	ChatsDispatcherStats() *dispatcher.ChatsDispatcherStats
	// ChatHandlerStats(chatID int64) *dispatcher.ApplicationChatHandlerStats
	MainLogs() string
	ChatLogs(chatID int64) string
	Shutdown()
}

type TgBotShutdowner interface {
	Shutdown()
}
