package emulator

import (
	"context"
	"math/rand"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	tgbot "github.com/nktknshn/go-tg-bot"
)

type FakeBot struct {
	lastID         int
	Messages       map[int]*models.Message
	updateCallback func()
}

func (fb *FakeBot) AddUserMessage(update *models.Update) {
	fb.Messages[update.Message.ID] = update.Message
}

// implement TelegramContextBot
func (fb *FakeBot) AnswerCallbackQuery(ctx context.Context, params *bot.AnswerCallbackQueryParams) (bool, error) {
	return true, nil
}

func (fb *FakeBot) SendMessage(ctx context.Context, params *bot.SendMessageParams) (*models.Message, error) {

	m := fb.createMessage(&tgbot.ChatRendererMessageProps{
		Text:        params.Text,
		ReplyMarkup: params.ReplyMarkup,
	})

	fb.notify()

	return m, nil
}

func tryInlineKeyboard(v models.ReplyMarkup) *models.InlineKeyboardMarkup {
	if e, ok := v.(models.InlineKeyboardMarkup); ok {
		return &e
	}
	return nil
}

func (fb *FakeBot) EditMessageText(ctx context.Context, params *bot.EditMessageTextParams) (*models.Message, error) {

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

func (fs *FakeBot) getNewID() int {
	fs.lastID++
	return fs.lastID
}

func (fs *FakeBot) createMessage(props *tgbot.ChatRendererMessageProps) *models.Message {

	botMessage := &models.Message{
		ID:          fs.getNewID(),
		Text:        props.Text,
		ReplyMarkup: *tryInlineKeyboard(props.ReplyMarkup),
	}

	fs.Messages[botMessage.ID] = botMessage

	return botMessage
}

func NewFakeBot() *FakeBot {
	return &FakeBot{
		Messages: make(map[int]*models.Message),
	}
}

type UpdateProps struct {
	ChatID int64
	UserID int64
}

type CallbackQueryUpdate struct {
	Data string
	UpdateProps
}

func NewCallbackQueryUpdate(props CallbackQueryUpdate) *models.Update {
	return &models.Update{
		ID: int64(rand.Int()),
		CallbackQuery: &models.CallbackQuery{
			Data: props.Data,
			Message: &models.Message{
				ID: rand.Int(),
				Chat: models.Chat{
					ID: props.ChatID,
				},
				From: &models.User{
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

// func NewTextMessageUpdateHelper(text string) *models.Update {

// }

func NewTextMessageUpdate(props TextMessageUpdate) *models.Update {
	return &models.Update{
		ID: int64(rand.Int()),
		Message: &models.Message{
			ID:   rand.Int(),
			Text: props.Text,
			Chat: models.Chat{
				ID: props.ChatID,
			},
			From: &models.User{
				ID:       int64(props.UserID),
				Username: "username",
			},
		},
	}
}
