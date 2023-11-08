package tgbot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var ErrMessageNotFound = fmt.Errorf("message not found")

type ChatRenderer interface {
	Message(context.Context, *ChatRendererMessageProps) (*models.Message, error)
	Delete(messageId int) error
}

type TelegramChatRenderer struct {
	Bot    *bot.Bot
	ChatID int64
}

type ChatRendererMessageProps struct {
	Text          string
	Extra         models.InlineKeyboardMarkup
	TargetMessage *models.Message
	RemoveTarget  bool
}

func (r *TelegramChatRenderer) Delete(messageId int) error {
	removed, err := r.Bot.DeleteMessage(context.Background(), &bot.DeleteMessageParams{
		ChatID:    r.ChatID,
		MessageID: messageId,
	})

	if err != nil {
		return fmt.Errorf("error removing target message %v: %w", messageId, err)
	}

	if !removed {
		return fmt.Errorf("error removing target message %v (false was returned)", messageId)
	}

	return nil
}

func (r *TelegramChatRenderer) Message(ctx context.Context, props *ChatRendererMessageProps) (*models.Message, error) {
	if props.TargetMessage != nil {

		// the message must be removed
		if props.RemoveTarget {
			err := r.Delete(props.TargetMessage.ID)

			if err != nil {
				return nil, err
			}
		} else {
			editedMessage, err := r.Bot.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:      r.ChatID,
				MessageID:   props.TargetMessage.ID,
				Text:        props.Text,
				ReplyMarkup: props.Extra,
			})

			if err != nil {
				return nil, fmt.Errorf("error editing message %v: %w", props.TargetMessage.ID, err)
			}

			return editedMessage, nil
		}

	}

	message, err := r.Bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      r.ChatID,
		Text:        props.Text,
		ReplyMarkup: props.Extra,
	})

	if err != nil {
		return nil, fmt.Errorf("error sending message: %w", err)
	}

	return message, nil
}
