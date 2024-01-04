package gotd

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"

	"github.com/nktknshn/go-tg-bot/helpers"
	"github.com/nktknshn/go-tg-bot/tgbot/dispatcher"
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
)

// Updates handler
type GotdHandler struct {
	dispatcher *dispatcher.ChatsDispatcher
	sender     *message.Sender
	client     *tg.Client
}

func NewGotdHandler(dispatcher *dispatcher.ChatsDispatcher) *GotdHandler {
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
