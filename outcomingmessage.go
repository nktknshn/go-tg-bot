package tgbot

import (
	"encoding/json"
	"fmt"

	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

// Messages that are going to be sent to user
// export type OutcomingMessageType = (OutcomingTextMessage<any> | OutcomingFileMessage) | OutcomingPhotoGroupMessage | OutcomingUserMessage

const (
	KindOutcomingFileMessage       = "OutcomingFileMessage"
	KindOutcomingUserMessage       = "OutcomingUserMessage"
	KindOutcomingPhotoGroupMessage = "OutcomingPhotoGroupMessage"
	KindOutcomingTextMessage       = "OutcomingTextMessage"
)

type OutcomingMessage interface {
	String() string
	OutcomingKind() string
	Equal(other OutcomingMessage) bool
}

type OutcomingFileMessage struct {
	ElementFile ElementFile
	Message     *models.Message
}

func (m OutcomingFileMessage) String() string {
	return fmt.Sprintf(
		"OutcomingFileMessage{ElementFile: %v, Message: %v}",
		m.ElementFile, m.Message,
	)
}

func (t *OutcomingFileMessage) OutcomingKind() string {
	return KindOutcomingFileMessage
}

func (t *OutcomingFileMessage) Equal(other OutcomingMessage) bool {
	if other.OutcomingKind() != KindOutcomingFileMessage {
		return false
	}

	otherFileMessage := other.(*OutcomingFileMessage)

	return t.ElementFile.FileId == otherFileMessage.ElementFile.FileId
}

type OutcomingUserMessage struct {
	ElementUserMessage ElementUserMessage
}

func (m OutcomingUserMessage) String() string {
	return fmt.Sprintf(
		"OutcomingUserMessage{ElementUserMessage: %v}",
		m.ElementUserMessage,
	)
}

func (t *OutcomingUserMessage) OutcomingKind() string {
	return KindOutcomingUserMessage
}

func (t *OutcomingUserMessage) Equal(other OutcomingMessage) bool {
	if other.OutcomingKind() != KindOutcomingUserMessage {
		return false
	}

	otherUserMessage := other.(*OutcomingUserMessage)

	return t.ElementUserMessage.MessageId == otherUserMessage.ElementUserMessage.MessageId
}

type OutcomingPhotoGroupMessage struct {
	ElementPhotoGroup ElementPhotoGroup
}

func (m OutcomingPhotoGroupMessage) String() string {
	return fmt.Sprintf(
		"OutcomingPhotoGroupMessage{ElementPhotoGroup: %v}",
		m.ElementPhotoGroup,
	)
}

func (t *OutcomingPhotoGroupMessage) OutcomingKind() string {
	return KindOutcomingPhotoGroupMessage
}

func (t *OutcomingPhotoGroupMessage) Equal(other OutcomingMessage) bool {
	if other.OutcomingKind() != KindOutcomingPhotoGroupMessage {
		return false
	}

	otherPhotoGroupMessage := other.(*OutcomingPhotoGroupMessage)

	return otherPhotoGroupMessage.ElementPhotoGroup.Equal(&t.ElementPhotoGroup)
}

type OutcomingTextMessage[A any] struct {
	Text       string
	Buttons    [][]*(ElementButton[A])
	isComplete bool
}

func (t *OutcomingTextMessage[T]) String() string {
	return fmt.Sprintf("OutcomingTextMessage{text: %s, buttons: %v, isComplete: %v}", t.Text, t.Buttons, t.isComplete)
}

func EqualInlineKeyboardMarkup(a models.InlineKeyboardMarkup, b models.InlineKeyboardMarkup) bool {
	if len(a.InlineKeyboard) != len(b.InlineKeyboard) {
		return false
	}

	for i, row := range a.InlineKeyboard {
		if len(row) != len(b.InlineKeyboard[i]) {
			return false
		}

		for j, button := range row {
			if button.Text != b.InlineKeyboard[i][j].Text {
				return false
			}

			if button.CallbackData != b.InlineKeyboard[i][j].CallbackData {
				return false
			}
		}
	}

	return true
}

func NewOutcomingTextMessage[A any](text string) *OutcomingTextMessage[A] {
	buttons := make([][]*ElementButton[A], 0)
	buttons = append(buttons, make([]*ElementButton[A], 0))

	return &OutcomingTextMessage[A]{
		Text:       text,
		Buttons:    buttons,
		isComplete: false,
	}
}

func (t *OutcomingTextMessage[T]) OutcomingKind() string {
	return KindOutcomingTextMessage
}

func (t *OutcomingTextMessage[T]) Equal(other OutcomingMessage) bool {
	if other.OutcomingKind() != KindOutcomingTextMessage {
		return false
	}

	oth := other.(*OutcomingTextMessage[T])

	return t.Text == oth.Text && EqualInlineKeyboardMarkup(t.getExtra(), oth.getExtra())
}

func (t *OutcomingTextMessage[T]) concatText(text string) {
	t.Text += "\n" + text
}

// func (t *OutcomingTextMessage[T]) complete() {
// 	t.isComplete = true
// }

func (t *OutcomingTextMessage[T]) getExtra() models.InlineKeyboardMarkup {
	res := models.InlineKeyboardMarkup{}

	for _, row := range t.Buttons {
		br := make([]models.InlineKeyboardButton, 0)
		for _, b := range row {
			br = append(br, models.InlineKeyboardButton{
				Text:         b.Text,
				CallbackData: b.CallbackData(),
			})
		}
		res.InlineKeyboard = append(res.InlineKeyboard, br)
	}

	return res
}

func (t *OutcomingTextMessage[T]) AddButton(button *ElementButton[T]) {
	t.Buttons[0] = append(t.Buttons[0], button)
}

func EqualExtra(a *models.ReplyMarkup, b *models.ReplyMarkup) bool {
	logger := GetLogger()

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	ab, err := json.Marshal(a)

	if err != nil {
		logger.Error("Error marshalling a", zap.Error(err))
		return false
	}

	bb, err := json.Marshal(b)

	if err != nil {
		logger.Error("Error marshalling b", zap.Error(err))
		return false
	}

	return string(ab) == string(bb)

}
