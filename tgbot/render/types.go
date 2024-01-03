package render

import (
	"context"

	"github.com/gotd/td/tg"
)

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
