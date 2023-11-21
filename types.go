package tgbot

import (
	"context"

	"github.com/go-telegram/bot/models"
)

// Interface for rendering messages into some interface (telegram, emulator, console, etc)
type ChatRenderer interface {
	Message(context.Context, *ChatRendererMessageProps) (*models.Message, error)
	Delete(messageId int) error
}

type ChatRendererMessageProps struct {
	Text          string
	ReplyMarkup   models.ReplyMarkup
	TargetMessage *models.Message
	RemoveTarget  bool
}

// Telegram entities
type RenderedElement interface {
	String() string
	renderedKind() string
	canReplace(outcomingMessage) bool
	Equal(RenderedElement) bool
	ID() int
}
