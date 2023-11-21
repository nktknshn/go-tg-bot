package tgbot

import (
	"fmt"

	"github.com/go-telegram/bot/models"
)

const (
	kindRenderedUserMessage        = "RenderedUserMessage"
	kindRenderedBotMessage         = "RenderedBotMessage"
	kindRenderedBotDocumentMessage = "RenderedBotDocumentMessage"
	kindRenderedPhotoGroup         = "RenderedPhotoGroup"
	// KindRenderedMediaMessage       = "RenderedMediaMessage"
)

// ta ne
// type R = RenderedUserMessage | RenderedBotMessage

type renderedUserMessage struct {
	OutcomingUserMessage *outcomingUserMessage
	MessageId            int
}

func newRenderedUserMessage(messageID int) *renderedUserMessage {
	return &renderedUserMessage{
		OutcomingUserMessage: &outcomingUserMessage{
			ElementUserMessage: elementUserMessage{
				MessageID: messageID,
			},
		},
		MessageId: messageID,
	}
}

func (r *renderedUserMessage) renderedKind() string {
	return kindRenderedUserMessage
}

func (r *renderedUserMessage) canReplace(outcoming outcomingMessage) bool {
	return false
}

func (r *renderedUserMessage) ID() int {
	return r.MessageId
}

// Equal
func (r *renderedUserMessage) Equal(other RenderedElement) bool {
	if other.renderedKind() != kindRenderedUserMessage {
		return false
	}
	otherUserMessage := other.(*renderedUserMessage)
	return r.MessageId == otherUserMessage.MessageId
}

func (r renderedUserMessage) String() string {
	return fmt.Sprintf(
		"RenderedUserMessage{MessageId: %v, OutcomingUserMessage: %v}",
		r.MessageId, r.OutcomingUserMessage,
	)
}

type renderedBotMessage struct {
	OutcomingTextMessage *outcomingTextMessage
	Message              *models.Message
}

func (rbm renderedBotMessage) String() string {
	return fmt.Sprintf(
		"RenderedBotMessage{OutcomingTextMessage: %v, Message: %v}",
		rbm.OutcomingTextMessage, rbm.Message.ID,
	)
}

func (r *renderedBotMessage) ID() int {
	return r.Message.ID
}

func (r *renderedBotMessage) renderedKind() string {
	return kindRenderedBotMessage
}

func (r *renderedBotMessage) canReplace(outcoming outcomingMessage) bool {
	return outcoming.OutcomingKind() == kindOutcomingTextMessage
	// && this.input.keyboardButtons.length == 0
	// && other.keyboardButtons.length == 0
	// TODO
}

// Equal
func (r *renderedBotMessage) Equal(other RenderedElement) bool {
	if other.renderedKind() != kindRenderedBotMessage {
		return false
	}
	otherBotMessage := other.(*renderedBotMessage)
	return r.OutcomingTextMessage.Equal(otherBotMessage.OutcomingTextMessage)
}

type renderedBotDocumentMessage struct {
	OutcomingFileMessage *outcomingFileMessage
}

func (r renderedBotDocumentMessage) String() string {
	return fmt.Sprintf(
		"RenderedBotDocumentMessage{OutcomingFileMessage: %v}",
		r.OutcomingFileMessage,
	)
}

func (r *renderedBotDocumentMessage) ID() int {
	return r.OutcomingFileMessage.Message.ID
}

func (r *renderedBotDocumentMessage) renderedKind() string {
	return kindRenderedBotDocumentMessage
}

func (r *renderedBotDocumentMessage) canReplace(outcoming outcomingMessage) bool {
	return false
}

// Equal
func (r *renderedBotDocumentMessage) Equal(other RenderedElement) bool {
	if other.renderedKind() != kindRenderedBotDocumentMessage {
		return false
	}
	otherBotDocumentMessage := other.(*renderedBotDocumentMessage)

	return r.OutcomingFileMessage.Equal(otherBotDocumentMessage.OutcomingFileMessage)
}

// type RenderedMediaMessage struct{}

// func (r *RenderedMediaMessage) renderedKind() string {
// 	return KindRenderedMediaMessage
// }

type renderedPhotoGroup struct {
	OutcomingPhotoGroupMessage *outcomingPhotoGroupMessage
	Message                    *models.Message
}

func (r *renderedPhotoGroup) ID() int {
	return r.Message.ID
}

func (r *renderedPhotoGroup) renderedKind() string {
	return kindRenderedPhotoGroup
}

func (r *renderedPhotoGroup) canReplace(outcoming outcomingMessage) bool {
	return false
}

// Equal TODO
func (r *renderedPhotoGroup) Equal(other RenderedElement) bool {
	if other.renderedKind() != kindRenderedPhotoGroup {
		return false
	}

	// otherPhotoGroup := other.(*RenderedPhotoGroup)
	return false
}
