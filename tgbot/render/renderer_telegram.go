package render

import (
	"context"
	"fmt"

	"github.com/gotd/td/tg"
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
)

// wrapping TelegramBot provides ChatRenderer interface
type telegramBotChatRenderer struct {
	Bot  telegram.TelegramBot
	User *tg.User
}

func (r *telegramBotChatRenderer) Delete(messageId int) error {
	removed, err := r.Bot.DeleteMessage(context.Background(), telegram.DeleteMessageParams{
		ChatID:     r.User.ID,
		AccessHash: r.User.AccessHash,
		MessageID:  messageId,
	})

	if err != nil {
		return fmt.Errorf("error removing target message %v: %w", messageId, err)
	}

	if !removed {
		return fmt.Errorf("error removing target message %v (false was returned)", messageId)
	}

	return nil
}

func (r *telegramBotChatRenderer) Message(ctx context.Context, props *ChatRendererMessageProps) (*tg.Message, error) {
	if props.TargetMessage != nil {

		// the message must be removed
		if props.RemoveTarget {
			err := r.Delete(props.TargetMessage.ID)

			if err != nil {
				return nil, err
			}
		} else {
			editedMessage, err := r.Bot.EditMessageText(ctx, telegram.EditMessageTextParams{
				ChatID:                r.User.ID,
				AccessHash:            r.User.AccessHash,
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

	message, err := r.Bot.SendMessage(ctx, telegram.SendMessageParams{
		ChatID:                r.User.ID,
		AccessHash:            r.User.AccessHash,
		Text:                  props.Text,
		ReplyMarkup:           props.ReplyMarkup,
		DisableWebPagePreview: true,
		// ParseMode:             tg.ParseModeMarkdown,
	})

	if err != nil {
		return nil, fmt.Errorf("error sending message: %w", err)
	}

	return message, nil
}

func NewTelegramChatRenderer(bot telegram.TelegramBot, user *tg.User) *telegramBotChatRenderer {
	return &telegramBotChatRenderer{
		Bot:  bot,
		User: user,
	}
}
