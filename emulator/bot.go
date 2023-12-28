package emulator

import (
	"context"
	"math/rand"

	"github.com/go-telegram/bot"

	tgbot "github.com/nktknshn/go-tg-bot"
)

type FakeBot struct {
	lastID         int
	Messages       map[int]*tg.Message
	updateCallback func()
	replyCallback  func()
	dispatcher     *tgbot.ChatsDispatcher
}

func (fb *FakeBot) Test3() {

}

func (fb *FakeBot) NewUser() *FakeBotUser {
	return &FakeBotUser{
		UserID: rand.Int63(),
		ChatID: rand.Int63(),
		Bot:    fb,
	}
}

func (fb *FakeBot) AddUserMessage(update *tg.Update) {
	fb.Messages[update.Message.ID] = update.Message
}

func (fb *FakeBot) SetDispatcher(d *tgbot.ChatsDispatcher) {
	fb.dispatcher = d
}

func (fb *FakeBot) AnswerCallbackQuery(ctx context.Context, params *bot.AnswerCallbackQueryParams) (bool, error) {

	if fb.replyCallback != nil {
		fb.replyCallback()
	}

	return true, nil
}

func (fb *FakeBot) SendMessage(ctx context.Context, params *bot.SendMessageParams) (*tg.Message, error) {

	m := fb.createMessage(&tgbot.ChatRendererMessageProps{
		Text:        params.Text,
		ReplyMarkup: params.ReplyMarkup,
	})

	fb.notify()

	return m, nil
}

func tryInlineKeyboard(v tg.ReplyMarkup) *tg.InlineKeyboardMarkup {
	if e, ok := v.(tg.InlineKeyboardMarkup); ok {
		return &e
	}
	return nil
}

func (fb *FakeBot) EditMessageText(ctx context.Context, params *bot.EditMessageTextParams) (*tg.Message, error) {

	if message, ok := fb.Messages[params.MessageID]; ok {
		message.Text = params.Text
		message.ReplyMarkup = *tryInlineKeyboard(params.ReplyMarkup)
		fb.notify()
		return message, nil
	}

	return nil, nil
}

func (fb *FakeBot) DeleteMessage(ctx context.Context, params *bot.DeleteMessageParams) (bool, error) {
	if _, ok := fb.Messages[params.MessageID]; ok {
		delete(fb.Messages, params.MessageID)
		fb.notify()
		return true, nil
	}

	return false, tgbot.ErrMessageNotFound
}

func (fb *FakeBot) notify() {
	if fb.updateCallback != nil {
		fb.updateCallback()
	}
}

func (fs *FakeBot) SetUpdateCallback(cb func()) {
	fs.updateCallback = cb
}

func (fs *FakeBot) SetReplyCallback(cb func()) {
	fs.replyCallback = cb
}

func (fs *FakeBot) getNewID() int {
	fs.lastID++
	return fs.lastID
}

func (fs *FakeBot) createMessage(props *tgbot.ChatRendererMessageProps) *tg.Message {

	botMessage := &tg.Message{
		ID:          fs.getNewID(),
		Text:        props.Text,
		ReplyMarkup: *tryInlineKeyboard(props.ReplyMarkup),
	}

	fs.Messages[botMessage.ID] = botMessage

	return botMessage
}

func NewFakeBot() *FakeBot {
	return &FakeBot{
		Messages: make(map[int]*tg.Message),
	}
}

type FakeBotUser struct {
	UserID int64
	ChatID int64
	Bot    *FakeBot
}

func (u *FakeBotUser) SendTextMessage(text string) *tg.Update {
	update := NewTextMessageUpdate(TextMessageUpdate{
		Text: text,
		UpdateProps: UpdateProps{
			ChatID: u.ChatID,
			UserID: u.UserID,
		},
	})

	u.Bot.AddUserMessage(update)
	u.Bot.dispatcher.HandleUpdate(context.Background(), u.Bot, update)

	return update
}

func (u *FakeBotUser) SendCallbackQuery(data string) *tg.Update {
	update := NewCallbackQueryUpdate(CallbackQueryUpdate{
		Data: data,
		UpdateProps: UpdateProps{
			ChatID: u.ChatID,
			UserID: u.UserID,
		},
	})

	u.Bot.dispatcher.HandleUpdate(context.Background(), u.Bot, update)

	return update
}

type UpdateProps struct {
	ChatID int64
	UserID int64
}

type CallbackQueryUpdate struct {
	Data string
	UpdateProps
}

func NewCallbackQueryUpdate(props CallbackQueryUpdate) *tg.Update {
	return &tg.Update{
		ID: int64(rand.Int()),
		CallbackQuery: &tg.CallbackQuery{
			Data: props.Data,
			Message: &tg.Message{
				ID: rand.Int(),
				Chat: tg.Chat{
					ID: props.ChatID,
				},
				From: &tg.User{
					ID:       int64(props.UserID),
					Username: "username",
				},
			},
		},
	}
}

type TextMessageUpdate struct {
	Text string
	UpdateProps
}

// func NewTextMessageUpdateHelper(text string) *tg.Update {

// }

func NewTextMessageUpdate(props TextMessageUpdate) *tg.Update {
	return &tg.Update{
		ID: int64(rand.Int()),
		Message: &tg.Message{
			ID:   rand.Int(),
			Text: props.Text,
			Chat: tg.Chat{
				ID: props.ChatID,
			},
			From: &tg.User{
				ID:       int64(props.UserID),
				Username: "username",
			},
		},
	}
}
