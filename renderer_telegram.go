package tgbot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type telegramChatRenderer struct {
	Bot    TelegramBot
	ChatID int64
}

func (r *telegramChatRenderer) Delete(messageId int) error {
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

func (r *telegramChatRenderer) Message(ctx context.Context, props *ChatRendererMessageProps) (*models.Message, error) {
	if props.TargetMessage != nil {

		// the message must be removed
		if props.RemoveTarget {
			err := r.Delete(props.TargetMessage.ID)

			if err != nil {
				return nil, err
			}
		} else {
			editedMessage, err := r.Bot.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:                r.ChatID,
				MessageID:             props.TargetMessage.ID,
				Text:                  props.Text,
				ReplyMarkup:           props.ReplyMarkup,
				DisableWebPagePreview: true,
			})

			if err != nil {
				return nil, fmt.Errorf("error editing message %v: %w", props.TargetMessage.ID, err)
			}

			return editedMessage, nil
		}

	}

	message, err := r.Bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:                r.ChatID,
		Text:                  props.Text,
		ReplyMarkup:           props.ReplyMarkup,
		DisableWebPagePreview: true,
		// ParseMode:             models.ParseModeMarkdown,
	})

	if err != nil {
		return nil, fmt.Errorf("error sending message: %w", err)
	}

	return message, nil
}

func NewTelegramChatRenderer(bot TelegramBot, chatID int64) *telegramChatRenderer {
	return &telegramChatRenderer{
		Bot:    bot,
		ChatID: chatID,
	}
}
