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
	kindOutcomingFileMessage       = "OutcomingFileMessage"
	kindOutcomingUserMessage       = "OutcomingUserMessage"
	kindOutcomingPhotoGroupMessage = "OutcomingPhotoGroupMessage"
	kindOutcomingTextMessage       = "OutcomingTextMessage"
)

type outcomingMessage interface {
	String() string
	OutcomingKind() string
	Equal(other outcomingMessage) bool
}

type outcomingFileMessage struct {
	ElementFile elementFile
	Message     *models.Message
}

func (m outcomingFileMessage) String() string {
	return fmt.Sprintf(
		"OutcomingFileMessage{ElementFile: %v, Message: %v}",
		m.ElementFile, m.Message,
	)
}

func (t *outcomingFileMessage) OutcomingKind() string {
	return kindOutcomingFileMessage
}

func (t *outcomingFileMessage) Equal(other outcomingMessage) bool {
	if other.OutcomingKind() != kindOutcomingFileMessage {
		return false
	}

	otherFileMessage := other.(*outcomingFileMessage)

	return t.ElementFile.FileId == otherFileMessage.ElementFile.FileId
}

type outcomingUserMessage struct {
	ElementUserMessage elementUserMessage
}

func (m outcomingUserMessage) String() string {
	return fmt.Sprintf(
		"OutcomingUserMessage{ElementUserMessage: %v}",
		m.ElementUserMessage,
	)
}

func (t *outcomingUserMessage) OutcomingKind() string {
	return kindOutcomingUserMessage
}

func (t *outcomingUserMessage) Equal(other outcomingMessage) bool {
	if other.OutcomingKind() != kindOutcomingUserMessage {
		return false
	}

	otherUserMessage := other.(*outcomingUserMessage)

	return t.ElementUserMessage.MessageID == otherUserMessage.ElementUserMessage.MessageID
}

type outcomingPhotoGroupMessage struct {
	ElementPhotoGroup elementPhotoGroup
}

func (m outcomingPhotoGroupMessage) String() string {
	return fmt.Sprintf(
		"OutcomingPhotoGroupMessage{ElementPhotoGroup: %v}",
		m.ElementPhotoGroup,
	)
}

func (t *outcomingPhotoGroupMessage) OutcomingKind() string {
	return kindOutcomingPhotoGroupMessage
}

func (t *outcomingPhotoGroupMessage) Equal(other outcomingMessage) bool {
	if other.OutcomingKind() != kindOutcomingPhotoGroupMessage {
		return false
	}

	otherPhotoGroupMessage := other.(*outcomingPhotoGroupMessage)

	return otherPhotoGroupMessage.ElementPhotoGroup.Equal(&t.ElementPhotoGroup)
}

type outcomingTextMessage struct {
	Text          string
	Buttons       [][]elementButton
	BottomButtons []elementBottomButton
	isComplete    bool
	// TODO RequestLocation
}

func (t *outcomingTextMessage) String() string {
	// maximum 20 chars from t.Text
	text := t.Text[:20]

	return fmt.Sprintf("OutcomingTextMessage{text: %s, buttons: %v, isComplete: %v}", text, t.Buttons, t.isComplete)
}

func equalReplyKeyboardMarkup(a models.ReplyKeyboardMarkup, b models.ReplyKeyboardMarkup) bool {
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

func equalInlineKeyboardMarkup(a models.InlineKeyboardMarkup, b models.InlineKeyboardMarkup) bool {
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

func newOutcomingTextMessage(text string) *outcomingTextMessage {
	buttons := make([][]elementButton, 0)
	buttons = append(buttons, make([]elementButton, 0))

	return &outcomingTextMessage{
		Text:       text,
		Buttons:    buttons,
		isComplete: false,
	}
}

func (t *outcomingTextMessage) OutcomingKind() string {
	return kindOutcomingTextMessage
}

func (t *outcomingTextMessage) Equal(other outcomingMessage) bool {
	if other.OutcomingKind() != kindOutcomingTextMessage {
		return false
	}

	oth := other.(*outcomingTextMessage)

	return t.Text == oth.Text && equalInlineKeyboardMarkup(t.InlineKeyboardMarkup(), oth.InlineKeyboardMarkup()) && equalReplyKeyboardMarkup(t.ReplyKeyboardMarkup(), oth.ReplyKeyboardMarkup())
}

func (t *outcomingTextMessage) ConcatText(text string) {
	t.Text += "\n" + text
}

func (t *outcomingTextMessage) SetComplete() {
	t.isComplete = true
}

func (t *outcomingTextMessage) ReplyMarkup() models.ReplyMarkup {

	if len(t.BottomButtons) > 0 {
		return t.ReplyKeyboardMarkup()
	}

	return t.InlineKeyboardMarkup()
}

func (t *outcomingTextMessage) ReplyKeyboardMarkup() models.ReplyKeyboardMarkup {
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

func (t *outcomingTextMessage) InlineKeyboardMarkup() models.InlineKeyboardMarkup {
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

func (t *outcomingTextMessage) AddButton(button *elementButton) {
	emptyButtons := len(t.Buttons) == 1 && len(t.Buttons[0]) == 0

	if button.NextRow && !emptyButtons {
		t.Buttons = append(t.Buttons, make([]elementButton, 0))
	}

	t.Buttons[len(t.Buttons)-1] = append(t.Buttons[len(t.Buttons)-1], *button)
}

func (t *outcomingTextMessage) AddButtonsRow(buttonsRow *elementButtonsRow) {
	t.Buttons = append(t.Buttons, buttonsRow.Buttons())
}

func equalReplyMarkup(a *models.ReplyMarkup, b *models.ReplyMarkup) bool {
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
