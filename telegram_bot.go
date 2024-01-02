package tgbot

import (
	"context"
	"crypto/rand"
	"io"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/unpack"
	"github.com/gotd/td/tg"
	"github.com/nktknshn/go-tg-bot/gogotd"
	"github.com/nktknshn/go-tg-bot/helpers"
)

// implements TelegramBot
type GotdBot struct {
	sender *message.Sender
	client *tg.Client
	rand   io.Reader
}

func NewGotdBot(sender *message.Sender, client *tg.Client) *GotdBot {
	return &GotdBot{sender: sender, client: client, rand: rand.Reader}
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

	outputMsg := &tg.MessagesEditMessageRequest{
		Peer:      &tg.InputPeerUser{UserID: params.ChatID, AccessHash: params.AccessHash},
		ID:        params.MessageID,
		Message:   params.Text,
		NoWebpage: params.DisableWebPagePreview,
	}

	println("params.ReplyMarkup", params.ReplyMarkup.TypeName())

	if !params.ReplyMarkup.Zero() {
		outputMsg.SetReplyMarkup(params.ReplyMarkup)
	}

	msg, err := gogotd.UnpackEditMessage(
		b.client.MessagesEditMessage(ctx, outputMsg),
	)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (b *GotdBot) SendMessage(ctx context.Context, params SendMessageParams) (*tg.Message, error) {

	id, err := RandInt64(b.rand)

	if err != nil {
		return nil, err
	}

	outcoming := &tg.MessagesSendMessageRequest{
		Peer:     &tg.InputPeerUser{UserID: params.ChatID, AccessHash: params.AccessHash},
		Message:  params.Text,
		RandomID: id,
	}

	if !params.ReplyMarkup.Zero() {
		outcoming.SetReplyMarkup(params.ReplyMarkup)
	}

	msg, err := unpack.Message(b.client.MessagesSendMessage(ctx, outcoming))

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
	return NewGotdBot(h.sender, h.client)
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

		h.dispatcher.HandleUpdate(ctx, h.Bot(), u)
		// multierr.AppendInto(&err, )
	}

	return err
}
