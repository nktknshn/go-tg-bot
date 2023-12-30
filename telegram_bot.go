package tgbot

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/unpack"
	"github.com/gotd/td/tg"
	"github.com/nktknshn/go-tg-bot/helpers"
	"go.uber.org/multierr"
)

// implements TelegramBot
type GotdBot struct {
	sender *message.Sender
	client *tg.Client
}

func (b *GotdBot) DeleteMessage(ctx context.Context, params DeleteMessageParams) (bool, error) {

	affected, err := b.client.MessagesDeleteMessages(ctx, &tg.MessagesDeleteMessagesRequest{
		ID:     []int{params.MessageID},
		Revoke: true,
	})

	if err != nil {
		return false, err
	}

	return affected.PtsCount > 0, nil
}

func (b *GotdBot) EditMessageText(ctx context.Context, params EditMessageTextParams) (*tg.Message, error) {

	msg, err := unpack.Message(
		b.client.MessagesEditMessage(ctx, &tg.MessagesEditMessageRequest{
			Peer:        &tg.InputPeerChat{ChatID: params.ChatID},
			ID:          params.MessageID,
			Message:     params.Text,
			ReplyMarkup: params.ReplyMarkup,
			NoWebpage:   params.DisableWebPagePreview,
		}),
	)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (b *GotdBot) SendMessage(ctx context.Context, params SendMessageParams) (*tg.Message, error) {

	msg, err := unpack.Message(b.client.MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
		Peer:        &tg.InputPeerChat{ChatID: params.ChatID},
		Message:     params.Text,
		ReplyMarkup: params.ReplyMarkup,
	}))

	if err != nil {
		return nil, err
	}

	return msg, nil
}

// reply to use button press
func (b *GotdBot) AnswerCallbackQuery(ctx context.Context, params AnswerCallbackQueryParams) (bool, error) {

	w, err := b.client.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
		QueryID: params.QueryID,
	})

	if err != nil {
		return false, err
	}

	return w, nil
}

// Updates handler
type Handler struct {
	dispatcher *ChatsDispatcher
	sender     *message.Sender
	client     *tg.Client
}

func (h *Handler) Bot() TelegramBot {
	return &GotdBot{sender: h.sender, client: h.client}
}

func (h *Handler) Handle(ctx context.Context, updates tg.UpdatesClass) error {

	extract, err := helpers.ExtractUpdates(updates)

	if err != nil {
		return err
	}

	for _, update := range extract.Updates {

		peerUser, ok := helpers.GetUser(extract.Entities, update)

		if !ok {
			continue
		}

		u := BotUpdate{
			UpdateClass: update,
			User:        peerUser,
			Entities:    extract.Entities,
		}

		multierr.AppendInto(&err, h.dispatcher.HandleUpdate(ctx, h.Bot(), u))
	}

	return err
}
