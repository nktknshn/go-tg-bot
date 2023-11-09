package tgbot

import (
	"fmt"

	"github.com/go-telegram/bot/models"
)

const (
	KindRenderedUserMessage        = "RenderedUserMessage"
	KindRenderedBotMessage         = "RenderedBotMessage"
	KindRenderedBotDocumentMessage = "RenderedBotDocumentMessage"
	KindRenderedPhotoGroup         = "RenderedPhotoGroup"
	// KindRenderedMediaMessage       = "RenderedMediaMessage"
)

// Telegram entities
type RenderedElement interface {
	String() string
	renderedKind() string
	canReplace(OutcomingMessage) bool
	Equal(RenderedElement) bool
	ID() int
}

// ta ne
// type R = RenderedUserMessage | RenderedBotMessage

type RenderedUserMessage struct {
	OutcomingUserMessage *OutcomingUserMessage
	MessageId            int
}

func NewRenderedUserMessage(messageID int) *RenderedUserMessage {
	return &RenderedUserMessage{
		OutcomingUserMessage: &OutcomingUserMessage{
			ElementUserMessage: ElementUserMessage{
				MessageID: messageID,
			},
		},
		MessageId: messageID,
	}
}

func (r *RenderedUserMessage) renderedKind() string {
	return KindRenderedUserMessage
}

func (r *RenderedUserMessage) canReplace(outcoming OutcomingMessage) bool {
	return false
}

func (r *RenderedUserMessage) ID() int {
	return r.MessageId
}

// Equal
func (r *RenderedUserMessage) Equal(other RenderedElement) bool {
	if other.renderedKind() != KindRenderedUserMessage {
		return false
	}
	otherUserMessage := other.(*RenderedUserMessage)
	return r.MessageId == otherUserMessage.MessageId
}

func (r RenderedUserMessage) String() string {
	return fmt.Sprintf(
		"RenderedUserMessage{MessageId: %v, OutcomingUserMessage: %v}",
		r.MessageId, r.OutcomingUserMessage,
	)
}

type RenderedBotMessage[A any] struct {
	OutcomingTextMessage *OutcomingTextMessage[A]
	Message              *models.Message
}

func (rbm RenderedBotMessage[A]) String() string {
	return fmt.Sprintf(
		"RenderedBotMessage{OutcomingTextMessage: %v, Message: %v}",
		rbm.OutcomingTextMessage, rbm.Message.ID,
	)
}

func (r *RenderedBotMessage[A]) ID() int {
	return r.Message.ID
}

func (r *RenderedBotMessage[A]) renderedKind() string {
	return KindRenderedBotMessage
}

func (r *RenderedBotMessage[A]) canReplace(outcoming OutcomingMessage) bool {
	return outcoming.OutcomingKind() == KindOutcomingTextMessage
	// && this.input.keyboardButtons.length == 0
	// && other.keyboardButtons.length == 0
	// TODO
}

// Equal
func (r *RenderedBotMessage[A]) Equal(other RenderedElement) bool {
	if other.renderedKind() != KindRenderedBotMessage {
		return false
	}
	otherBotMessage := other.(*RenderedBotMessage[A])
	return r.OutcomingTextMessage.Equal(otherBotMessage.OutcomingTextMessage)
}

type RenderedBotDocumentMessage struct {
	OutcomingFileMessage *OutcomingFileMessage
}

func (r RenderedBotDocumentMessage) String() string {
	return fmt.Sprintf(
		"RenderedBotDocumentMessage{OutcomingFileMessage: %v}",
		r.OutcomingFileMessage,
	)
}

func (r *RenderedBotDocumentMessage) ID() int {
	return r.OutcomingFileMessage.Message.ID
}

func (r *RenderedBotDocumentMessage) renderedKind() string {
	return KindRenderedBotDocumentMessage
}

func (r *RenderedBotDocumentMessage) canReplace(outcoming OutcomingMessage) bool {
	return false
}

// Equal
func (r *RenderedBotDocumentMessage) Equal(other RenderedElement) bool {
	if other.renderedKind() != KindRenderedBotDocumentMessage {
		return false
	}
	otherBotDocumentMessage := other.(*RenderedBotDocumentMessage)

	return r.OutcomingFileMessage.Equal(otherBotDocumentMessage.OutcomingFileMessage)
}

// type RenderedMediaMessage struct{}

// func (r *RenderedMediaMessage) renderedKind() string {
// 	return KindRenderedMediaMessage
// }

type RenderedPhotoGroup struct {
	OutcomingPhotoGroupMessage *OutcomingPhotoGroupMessage
	Message                    *models.Message
}

func (r *RenderedPhotoGroup) ID() int {
	return r.Message.ID
}

func (r *RenderedPhotoGroup) renderedKind() string {
	return KindRenderedPhotoGroup
}

func (r *RenderedPhotoGroup) canReplace(outcoming OutcomingMessage) bool {
	return false
}

// Equal TODO
func (r *RenderedPhotoGroup) Equal(other RenderedElement) bool {
	if other.renderedKind() != KindRenderedPhotoGroup {
		return false
	}

	// otherPhotoGroup := other.(*RenderedPhotoGroup)
	return false
}
