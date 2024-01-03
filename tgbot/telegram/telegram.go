package telegram

import (
	"context"

	"github.com/gotd/td/tg"
)

type AnswerCallbackQueryParams struct {
	QueryID int64
}

type CallbackAnswerer interface {
	AnswerCallbackQuery(context.Context, AnswerCallbackQueryParams) (bool, error)
}

// Interface for rendering messages into some interface (telegram, emulator, console, etc)
type TelegramBot interface {
	MessageDeleter
	MessageEditor
	MessageSender
	CallbackAnswerer
}

type TelegramUserChat struct {
	ChatID int64
	Chat   *tg.Chat
}
