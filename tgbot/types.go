package tgbot

import (
	"context"

	"github.com/gotd/td/tg"
)

type BotUpdate struct {
	UpdateClass tg.UpdateClass
	User        *tg.User
	Entities    tg.Entities
}

func (bu BotUpdate) GetNewMessageUpdate() (*tg.UpdateNewMessage, bool) {
	if update, ok := bu.UpdateClass.(*tg.UpdateNewMessage); ok {
		return update, true
	}

	return nil, false
}

func (bu BotUpdate) GetCallbackQueryUpdate() (*tg.UpdateBotCallbackQuery, bool) {
	if update, ok := bu.UpdateClass.(*tg.UpdateBotCallbackQuery); ok {
		return update, true
	}

	return nil, false
}

// Interface for rendering messages into some interface (telegram, emulator, console, etc)
type ChatRenderer interface {
	Message(context.Context, *ChatRendererMessageProps) (*tg.Message, error)
	Delete(messageId int) error
}

type ChatRendererMessageProps struct {
	Text          string
	ReplyMarkup   tg.ReplyMarkupClass
	TargetMessage *tg.Message
	RemoveTarget  bool
}

// Telegram entities
