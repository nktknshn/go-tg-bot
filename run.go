package tgbot

import (
	"context"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/unpack"
	"github.com/gotd/td/tg"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type HandlerBot struct {
	sender *message.Sender
	client *tg.Client
}

func (b *HandlerBot) DeleteMessage(ctx context.Context, params DeleteMessageParams) (bool, error) {

	affected, err := b.client.MessagesDeleteMessages(ctx, &tg.MessagesDeleteMessagesRequest{
		ID:     []int{params.MessageID},
		Revoke: true,
	})

	if err != nil {
		return false, err
	}

	return affected.PtsCount > 0, nil
}

func (b *HandlerBot) EditMessageText(ctx context.Context, params EditMessageTextParams) (*tg.Message, error) {

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

func (b *HandlerBot) SendMessage(ctx context.Context, params SendMessageParams) (*tg.Message, error) {

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
func (b *HandlerBot) AnswerCallbackQuery(ctx context.Context, params AnswerCallbackQueryParams) (bool, error) {

	w, err := b.client.MessagesSetBotCallbackAnswer(ctx, &tg.MessagesSetBotCallbackAnswerRequest{
		QueryID: params.QueryID,
	})

	if err != nil {
		return false, err
	}

	return w, nil
}

type Handler struct {
	dispatcher *ChatsDispatcher
	sender     *message.Sender
	client     *tg.Client
}

func (h *Handler) Bot() TelegramBot {
	return &HandlerBot{sender: h.sender, client: h.client}
}

func (h *Handler) Handle(ctx context.Context, updates tg.UpdatesClass) error {

	var (
		e    tg.Entities
		upds []tg.UpdateClass
	)

	switch u := updates.(type) {
	case *tg.Updates:
		upds = u.Updates
		e.Users = u.MapUsers().NotEmptyToMap()
		chats := u.MapChats()
		e.Chats = chats.ChatToMap()
		e.Channels = chats.ChannelToMap()
	case *tg.UpdatesCombined:
		upds = u.Updates
		e.Users = u.MapUsers().NotEmptyToMap()
		chats := u.MapChats()
		e.Chats = chats.ChatToMap()
		e.Channels = chats.ChannelToMap()
	case *tg.UpdateShort:
		upds = []tg.UpdateClass{u.Update}
	default:
		// *UpdateShortMessage
		// *UpdateShortChatMessage
		// *UpdateShortSentMessage
		// *UpdatesTooLong
		return nil
	}

	var err error
	for _, update := range upds {
		multierr.AppendInto(&err, h.dispatcher.HandleUpdate(ctx, h.Bot(), update))
	}

	return err
}

func Run(ctx context.Context, logger *zap.Logger, dispatcher *ChatsDispatcher) error {
	logger.Debug("Starting real telegram bot")

	handler := &Handler{
		dispatcher: dispatcher,
	}

	opts := telegram.Options{
		Logger:        logger,
		UpdateHandler: handler,
	}

	err := telegram.BotFromEnvironment(ctx, opts,
		func(ctx context.Context, client *telegram.Client) error {

			api := tg.NewClient(client)

			sender := message.NewSender(api)
			handler.sender = sender
			handler.client = api

			return nil
		},
		func(ctx context.Context, client *telegram.Client) error {
			return nil
		},
	)

	return err

	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	// defer cancel()

	// token := os.Getenv("BOT_TOKEN")

	// if token == "" {
	// 	logger.Fatal("BOT_TOKEN env variable is not set")
	// 	os.Exit(1)
	// }

	// bot, err := bot.New(token, bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *tg.Update) {
	// 	dispatcher.HandleUpdate(ctx, bot, update)
	// }))

	// if err != nil {
	// 	logger.Fatal("Error creating bot", zap.Error(err))
	// 	os.Exit(1)
	// }

	// bot.Start(ctx)

}
