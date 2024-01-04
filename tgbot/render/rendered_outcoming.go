package render

import (
	"fmt"

	"github.com/gotd/td/tg"
	"github.com/nktknshn/go-tg-bot/tgbot/component"
	"github.com/nktknshn/go-tg-bot/tgbot/outcoming"
)

const (
	KindRenderedUserMessage        = "RenderedUserMessage"
	KindRenderedBotMessage         = "RenderedBotMessage"
	KindRenderedBotDocumentMessage = "RenderedBotDocumentMessage"
	KindRenderedPhotoGroup         = "RenderedPhotoGroup"
	// KindRenderedMediaMessage       = "RenderedMediaMessage"
)

type RenderedElement interface {
	ID() int
	String() string
	RenderedKind() string
	CanReplace(outcoming.OutcomingMessage) bool
	Equal(RenderedElement) bool
}

// Messages that have been sent to user

type RenderedUserMessage struct {
	OutcomingUserMessage *outcoming.OutcomingUserMessage
	MessageID            int
}

func NewRenderedUserMessage(messageID int) *RenderedUserMessage {
	return &RenderedUserMessage{
		OutcomingUserMessage: &outcoming.OutcomingUserMessage{
			ElementUserMessage: component.ElementUserMessage{
				MessageID: messageID,
			},
		},
		MessageID: messageID,
	}
}

func (r *RenderedUserMessage) RenderedKind() string {
	return KindRenderedUserMessage
}

func (r *RenderedUserMessage) CanReplace(outcoming outcoming.OutcomingMessage) bool {
	return false
}

func (r *RenderedUserMessage) ID() int {
	return r.MessageID
}

// Equal
func (r *RenderedUserMessage) Equal(other RenderedElement) bool {
	if other.RenderedKind() != KindRenderedUserMessage {
		return false
	}
	otherUserMessage := other.(*RenderedUserMessage)
	return r.MessageID == otherUserMessage.MessageID
}

func (r RenderedUserMessage) String() string {
	return fmt.Sprintf(
		"RenderedUserMessage{MessageId: %v, OutcomingUserMessage: %v}",
		r.MessageID, r.OutcomingUserMessage,
	)
}

type RenderedBotMessage struct {
	OutcomingTextMessage *outcoming.OutcomingTextMessage
	Message              *tg.Message
}

func (rbm RenderedBotMessage) String() string {
	return fmt.Sprintf(
		"RenderedBotMessage{OutcomingTextMessage: %v, Message: %v}",
		rbm.OutcomingTextMessage, rbm.Message.ID,
	)
}

func (r *RenderedBotMessage) ID() int {
	return r.Message.ID
}

func (r *RenderedBotMessage) RenderedKind() string {
	return KindRenderedBotMessage
}

func (r *RenderedBotMessage) CanReplace(out outcoming.OutcomingMessage) bool {
	return out.OutcomingKind() == outcoming.KindOutcomingTextMessage
	// && this.input.keyboardButtons.length == 0
	// && other.keyboardButtons.length == 0
	// TODO
}

// Equal
func (r *RenderedBotMessage) Equal(other RenderedElement) bool {
	if other.RenderedKind() != KindRenderedBotMessage {
		return false
	}
	otherBotMessage := other.(*RenderedBotMessage)
	return r.OutcomingTextMessage.Equal(otherBotMessage.OutcomingTextMessage)
}

type RenderedBotDocumentMessage struct {
	OutcomingFileMessage *outcoming.OutcomingFileMessage
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

func (r *RenderedBotDocumentMessage) RenderedKind() string {
	return KindRenderedBotDocumentMessage
}

func (r *RenderedBotDocumentMessage) CanReplace(outcoming outcoming.OutcomingMessage) bool {
	return false
}

// Equal
func (r *RenderedBotDocumentMessage) Equal(other RenderedElement) bool {
	if other.RenderedKind() != KindRenderedBotDocumentMessage {
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
	OutcomingPhotoGroupMessage *outcoming.OutcomingPhotoGroupMessage
	Message                    *tg.Message
}

func (r *RenderedPhotoGroup) ID() int {
	return r.Message.ID
}

func (r *RenderedPhotoGroup) renderedKind() string {
	return KindRenderedPhotoGroup
}

func (r *RenderedPhotoGroup) CanReplace(outcoming outcoming.OutcomingMessage) bool {
	return false
}

// Equal TODO
func (r *RenderedPhotoGroup) Equal(other RenderedElement) bool {
	if other.RenderedKind() != KindRenderedPhotoGroup {
		return false
	}

	// otherPhotoGroup := other.(*RenderedPhotoGroup)
	return false
}
