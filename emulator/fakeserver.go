package emulator

import (
	"context"

	"github.com/go-telegram/bot/models"
	tgbot "github.com/nktknshn/go-tg-bot"
	"go.uber.org/zap"
)

/*
This is a fake telegram server for testing purposes
*/

type FakeServer struct {
	lastID         int
	Messages       map[int]*models.Message
	updateCallback func()
}

func NewFakeServer() *FakeServer {
	return &FakeServer{
		lastID:   0,
		Messages: make(map[int]*models.Message),
	}
}

func (fs *FakeServer) SetUpdateCallback(cb func()) {
	fs.updateCallback = cb
}

func (fs *FakeServer) getNewID() int {
	fs.lastID++
	return fs.lastID
}

func (fs *FakeServer) createMessage(props *tgbot.ChatRendererMessageProps) *models.Message {
	botMessage := &models.Message{
		ID:          fs.getNewID(),
		Text:        props.Text,
		ReplyMarkup: props.Extra,
	}

	fs.Messages[botMessage.ID] = botMessage

	return botMessage
}

func (fs *FakeServer) Delete(messageID int) error {
	if _, ok := fs.Messages[messageID]; ok {
		delete(fs.Messages, messageID)

		if fs.updateCallback != nil {
			fs.updateCallback()
		}

		return nil
	}

	return tgbot.ErrMessageNotFound
}

func (fs *FakeServer) Message(ctx context.Context, props *tgbot.ChatRendererMessageProps) (*models.Message, error) {

	var editMessage bool
	var deleteMessage bool
	var m *models.Message

	logger.Debug("FakeServer.Message", zap.Any("props", props), zap.Int("total", len(fs.Messages)))

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
			m = fs.createMessage(props)
		} else if editMessage {
			m = fs.Messages[targetID]
			m.Text = props.Text
		}
	} else {
		m = fs.createMessage(props)
	}

	if fs.updateCallback != nil {
		fs.updateCallback()
	}

	logger.Debug("Returning", zap.Any("m", m))

	return m, nil
}
