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
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
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

func (b *GotdBot) DeleteMessage(ctx context.Context, params telegram.DeleteMessageParams) (bool, error) {

	affected, err := b.client.MessagesDeleteMessages(ctx, &tg.MessagesDeleteMessagesRequest{
		ID:     []int{params.MessageID},
		Revoke: true,
	})

	if err != nil {
		return false, err
	}

	return affected.PtsCount > 0, nil
}

func (b *GotdBot) EditMessageText(ctx context.Context, params telegram.EditMessageTextParams) (*tg.Message, error) {

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

func (b *GotdBot) SendMessage(ctx context.Context, params telegram.SendMessageParams) (*tg.Message, error) {

	id, err := helpers.RandInt64(b.rand)

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
func (b *GotdBot) AnswerCallbackQuery(ctx context.Context, params telegram.AnswerCallbackQueryParams) (bool, error) {

	w, err := b.client.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
		QueryID: params.QueryID,
	})

	if err != nil {
		return false, err
	}

	return w, nil
}

// Updates handler
type GotdHandler struct {
	dispatcher *ChatsDispatcher
	sender     *message.Sender
	client     *tg.Client
}

func NewGotdHandler(dispatcher *ChatsDispatcher) *GotdHandler {
	return &GotdHandler{dispatcher: dispatcher}
}

func (h *GotdHandler) SetSender(sender *message.Sender) {
	h.sender = sender
}

func (h *GotdHandler) SetClient(client *tg.Client) {
	h.client = client
}

func (h *GotdHandler) Bot() telegram.TelegramBot {
	return NewGotdBot(h.sender, h.client)
}

func (h *GotdHandler) Handle(ctx context.Context, updates tg.UpdatesClass) error {

	extract, err := helpers.ExtractUpdates(updates)

	if err != nil {
		return err
	}

	for _, update := range extract.Updates {

		peerUser, ok := helpers.GetUser(extract.Entities, update)

		if !ok {
			continue
		}

		u := telegram.BotUpdate{
			UpdateClass: update,
			User:        peerUser,
			Entities:    extract.Entities,
		}

		h.dispatcher.HandleUpdate(ctx, h.Bot(), u)
		// multierr.AppendInto(&err, )
	}

	return err
}
