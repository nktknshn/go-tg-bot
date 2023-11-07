package emulator

import (
	"context"

	"github.com/go-telegram/bot/models"
	tgbot "github.com/nktknshn/go-tg-bot"
)

/*
This is a fake telegram server for testing purposes
*/

type FakeServer struct {
	lastID   int
	Messages map[int]*models.Message
}

func NewFakeServer() *FakeServer {
	return &FakeServer{}
}

func (fs *FakeServer) getNewID() int {
	fs.lastID++
	return fs.lastID
}

func (fs *FakeServer) createMessage(props *tgbot.ChatRendererMessageProps) *models.Message {
	botMessage := &models.Message{
		ID:          fs.getNewID(),
		Text:        props.Text,
		ReplyMarkup: models.InlineKeyboardMarkup{},
	}

	fs.Messages[botMessage.ID] = botMessage

	return botMessage
}

func (fs *FakeServer) Delete(messageID int) error {
	if _, ok := fs.Messages[messageID]; ok {
		delete(fs.Messages, messageID)
		return nil
	}

	return tgbot.ErrMessageNotFound
}

func (fs *FakeServer) Message(ctx context.Context, props *tgbot.ChatRendererMessageProps) (*models.Message, error) {

	var editMessage bool
	var deleteMessage bool

	if props.TargetMessage != nil {
		targetID := props.TargetMessage.ID

		if _, ok := fs.Messages[targetID]; ok {
			// delete(fs.Messages, targetID)
			if props.RemoveTarget {
				deleteMessage = true
			} else {
				editMessage = true
			}
		} else {
			return nil, tgbot.ErrMessageNotFound
		}

		if deleteMessage {
			delete(fs.Messages, targetID)

		} else if editMessage {
			m := fs.Messages[targetID]
			m.Text = props.Text
		}
	}

	return fs.createMessage(props), nil
}
