package telegram

import (
	"context"
	"fmt"

	"github.com/gotd/td/tg"
)

var ErrMessageNotFound = fmt.Errorf("message not found")

type DeleteMessageParams struct {
	ChatID     int64
	AccessHash int64
	MessageID  int
}

type EditMessageTextParams struct {
	ChatID     int64
	AccessHash int64

	MessageID             int
	Text                  string
	ReplyMarkup           tg.ReplyMarkupClass
	DisableWebPagePreview bool
}

type SendMessageParams struct {
	ChatID     int64
	AccessHash int64

	Text                  string
	ReplyMarkup           tg.ReplyMarkupClass
	DisableWebPagePreview bool
}

type MessageDeleter interface {
	DeleteMessage(ctx context.Context, params DeleteMessageParams) (bool, error)
}

type MessageEditor interface {
	EditMessageText(ctx context.Context, params EditMessageTextParams) (*tg.Message, error)
}

type MessageSender interface {
	SendMessage(ctx context.Context, params SendMessageParams) (*tg.Message, error)
}
