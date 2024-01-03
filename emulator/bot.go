package emulator

import (
	"context"

	"github.com/gotd/td/tg"

	tgbot "github.com/nktknshn/go-tg-bot/tgbot"
)

type FakeBot struct {
	dispatcher *tgbot.ChatsDispatcher

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

func (fb *FakeBot) SetDispatcher(d *tgbot.ChatsDispatcher) {
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

func (fb *FakeBot) SendMessage(ctx context.Context, params tgbot.SendMessageParams) (*tg.Message, error) {

	m := fb.createMessage(params.ChatID, &tgbot.ChatRendererMessageProps{
		Text:        params.Text,
		ReplyMarkup: params.ReplyMarkup,
	})

	fb.notify()

	return m, nil
}

func (fb *FakeBot) AnswerCallbackQuery(ctx context.Context, params tgbot.AnswerCallbackQueryParams) (bool, error) {

	if fb.replyCallback != nil {
		fb.replyCallback()
	}

	return true, nil
}

func (fb *FakeBot) EditMessageText(ctx context.Context, params tgbot.EditMessageTextParams) (*tg.Message, error) {

	if message, ok := fb.Messages[params.MessageID]; ok {
		message.Message = params.Text
		message.ReplyMarkup = params.ReplyMarkup

		fb.notify()
		return message, nil
	}

	return nil, nil
}

func (fb *FakeBot) DeleteMessage(ctx context.Context, params tgbot.DeleteMessageParams) (bool, error) {
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
	fs.lastMessageID++
	return fs.lastMessageID
}

func (fs *FakeBot) createMessage(chatID int64, props *tgbot.ChatRendererMessageProps) *tg.Message {
	botMessage := &tg.Message{
		ID:          fs.getNewID(),
		Message:     props.Text,
		ReplyMarkup: props.ReplyMarkup,
		PeerID:      &tg.PeerUser{UserID: chatID},
	}

	fs.Messages[botMessage.ID] = botMessage

	return botMessage
}
