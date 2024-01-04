package emulator

import (
	"context"

	"github.com/gotd/td/tg"

	"github.com/nktknshn/go-tg-bot/tgbot/dispatcher"
	"github.com/nktknshn/go-tg-bot/tgbot/render"
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
)

type FakeBot struct {
	dispatcher *dispatcher.ChatsDispatcher

	lastMessageID int

	Messages       map[int]*tg.Message
	updateCallback func()
	replyCallback  func()
}

func NewFakeBot() *FakeBot {
	return &FakeBot{
		Messages: make(map[int]*tg.Message),
	}
}

func (fb *FakeBot) NewUser(userID int64) *FakeBotUser {
	u := NewFakeBotUser(userID)
	u.Bot = fb
	return u
}

func (fb *FakeBot) AddUserMessage(message *tg.Message) {
	fb.Messages[message.ID] = &tg.Message{
		ID:      fb.getNewID(),
		Message: message.Message,
		PeerID:  message.PeerID,
	}
}

func (fb *FakeBot) SetDispatcher(d *dispatcher.ChatsDispatcher) {
	fb.dispatcher = d
}

func (fb *FakeBot) DisplayedMessages(chatID int64) []*tg.Message {
	var messages []*tg.Message

	for _, message := range fb.Messages {
		if message.PeerID.(*tg.PeerUser).UserID == chatID {
			messages = append(messages, message)
		}
	}

	return messages
}

func (fb *FakeBot) SendMessage(ctx context.Context, params telegram.SendMessageParams) (*tg.Message, error) {

	m := fb.createMessage(params.ChatID, &render.ChatRendererMessageProps{
		Text:        params.Text,
		ReplyMarkup: params.ReplyMarkup,
	})

	fb.notify()

	return m, nil
}

func (fb *FakeBot) AnswerCallbackQuery(ctx context.Context, params telegram.AnswerCallbackQueryParams) (bool, error) {

	if fb.replyCallback != nil {
		fb.replyCallback()
	}

	return true, nil
}

func (fb *FakeBot) EditMessageText(ctx context.Context, params telegram.EditMessageTextParams) (*tg.Message, error) {

	if message, ok := fb.Messages[params.MessageID]; ok {
		message.Message = params.Text
		message.ReplyMarkup = params.ReplyMarkup

		fb.notify()
		return message, nil
	}

	return nil, nil
}

func (fb *FakeBot) DeleteMessage(ctx context.Context, params telegram.DeleteMessageParams) (bool, error) {
	if _, ok := fb.Messages[params.MessageID]; ok {
		delete(fb.Messages, params.MessageID)
		fb.notify()
		return true, nil
	}

	return false, telegram.ErrMessageNotFound
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
	fs.lastMessageID++
	return fs.lastMessageID
}

func (fs *FakeBot) createMessage(chatID int64, props *render.ChatRendererMessageProps) *tg.Message {
	botMessage := &tg.Message{
		ID:          fs.getNewID(),
		Message:     props.Text,
		ReplyMarkup: props.ReplyMarkup,
		PeerID:      &tg.PeerUser{UserID: chatID},
	}

	fs.Messages[botMessage.ID] = botMessage

	return botMessage
}
