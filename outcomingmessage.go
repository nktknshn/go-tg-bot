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

	return t.ElementUserMessage.MessageID == otherUserMessage.ElementUserMessage.MessageID
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
	Text          string
	Buttons       [][]ElementButton[A]
	BottomButtons []ElementBottomButton
	isComplete    bool
	// TODO RequestLocation
}

func (t *OutcomingTextMessage[T]) String() string {
	// maximum 20 chars from t.Text
	text := t.Text[:20]

	return fmt.Sprintf("OutcomingTextMessage{text: %s, buttons: %v, isComplete: %v}", text, t.Buttons, t.isComplete)
}

func EqualReplyKeyboardMarkup(a models.ReplyKeyboardMarkup, b models.ReplyKeyboardMarkup) bool {
	if len(a.Keyboard) != len(b.Keyboard) {
		return false
	}

	for i, row := range a.Keyboard {
		if len(row) != len(b.Keyboard[i]) {
			return false
		}

		for j, button := range row {
			if button.Text != b.Keyboard[i][j].Text {
				return false
			}
		}
	}

	return true
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
	buttons := make([][]ElementButton[A], 0)
	buttons = append(buttons, make([]ElementButton[A], 0))

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

	return t.Text == oth.Text && EqualInlineKeyboardMarkup(t.InlineKeyboardMarkup(), oth.InlineKeyboardMarkup()) && EqualReplyKeyboardMarkup(t.ReplyKeyboardMarkup(), oth.ReplyKeyboardMarkup())
}

func (t *OutcomingTextMessage[T]) ConcatText(text string) {
	t.Text += "\n" + text
}

func (t *OutcomingTextMessage[T]) SetComplete() {
	t.isComplete = true
}

func (t *OutcomingTextMessage[T]) ReplyMarkup() models.ReplyMarkup {

	if len(t.BottomButtons) > 0 {
		return t.ReplyKeyboardMarkup()
	}

	return t.InlineKeyboardMarkup()
}

func (t *OutcomingTextMessage[T]) ReplyKeyboardMarkup() models.ReplyKeyboardMarkup {
	res := models.ReplyKeyboardMarkup{}

	if len(t.BottomButtons) > 0 {

		for _, b := range t.BottomButtons {
			if len(b.Texts) > 0 {
				br := make([]models.KeyboardButton, 0)
				for _, t := range b.Texts {
					br = append(br, models.KeyboardButton{Text: t})
				}
				res.Keyboard = append(res.Keyboard, br)
			} else {
				res.Keyboard = append(res.Keyboard, []models.KeyboardButton{{Text: b.Text}})
			}
		}

	}
	return res
}

func (t *OutcomingTextMessage[T]) InlineKeyboardMarkup() models.InlineKeyboardMarkup {
	// ReplyKeyboardRemove
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
	emptyButtons := len(t.Buttons) == 1 && len(t.Buttons[0]) == 0

	if button.NextRow && !emptyButtons {
		t.Buttons = append(t.Buttons, make([]ElementButton[T], 0))
	}

	t.Buttons[len(t.Buttons)-1] = append(t.Buttons[len(t.Buttons)-1], *button)
}

func (t *OutcomingTextMessage[T]) AddButtonsRow(buttonsRow *ElementButtonsRow[T]) {
	t.Buttons = append(t.Buttons, buttonsRow.Buttons())
}

func EqualReplyMarkup(a *models.ReplyMarkup, b *models.ReplyMarkup) bool {
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
