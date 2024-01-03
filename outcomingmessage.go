package tgbot

import (
	"encoding/json"
	"fmt"

	// "github.com/gotd/td/tg"
	"github.com/gotd/td/bin"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

// Messages that are going to be sent to user
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
	Message     *tg.Message
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
	text := t.Text[:min(20, len(t.Text))]

	return fmt.Sprintf("OutcomingTextMessage{text: %s, buttons: %v, isComplete: %v}", text, t.Buttons, t.isComplete)
}

func equalReplyKeyboardMarkup(a *tg.ReplyKeyboardMarkup, b *tg.ReplyKeyboardMarkup) bool {
	if len(a.Rows) != len(b.Rows) {
		return false
	}

	for i, row := range a.Rows {
		if len(row.Buttons) != len(b.Rows[i].Buttons) {
			return false
		}

		for j, button := range row.Buttons {
			if button.GetText() != b.Rows[i].Buttons[j].GetText() {
				return false
			}
		}
	}

	return true
}

func equalInlineKeyboardButtonClass(a tg.KeyboardButtonClass, b tg.KeyboardButtonClass) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	b1 := bin.Buffer{
		Buf: make([]byte, 0),
	}

	b2 := bin.Buffer{
		Buf: make([]byte, 0),
	}

	err := a.Encode(&b1)

	if err != nil {
		panic(err)
	}

	err = b.Encode(&b2)

	if err != nil {
		panic(err)
	}

	return string(b1.Buf) == string(b2.Buf)
}

func equalInlineKeyboardMarkup(a *tg.ReplyInlineMarkup, b *tg.ReplyInlineMarkup) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a.Rows) != len(b.Rows) {
		return false
	}

	for i, row := range a.Rows {
		if len(row.Buttons) != len(b.Rows[i].Buttons) {
			return false
		}

		for j, button := range row.Buttons {
			if button.GetText() != b.Rows[i].Buttons[j].GetText() {
				return false
			}

			if !equalInlineKeyboardButtonClass(button, b.Rows[i].Buttons[j]) {
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

	return t.Text == oth.Text &&
		equalInlineKeyboardMarkup(t.InlineKeyboardMarkup(), oth.InlineKeyboardMarkup()) &&
		equalReplyKeyboardMarkup(t.ReplyKeyboardMarkup(), oth.ReplyKeyboardMarkup())
}

func (t *outcomingTextMessage) ConcatText(text string) {
	t.Text += "\n" + text
}

func (t *outcomingTextMessage) SetComplete() {
	t.isComplete = true
}

func (t *outcomingTextMessage) ReplyMarkup() tg.ReplyMarkupClass {

	if len(t.BottomButtons) > 0 {
		return t.ReplyKeyboardMarkup()
	}

	return t.InlineKeyboardMarkup()
}

func (t *outcomingTextMessage) ReplyKeyboardMarkup() *tg.ReplyKeyboardMarkup {
	res := tg.ReplyKeyboardMarkup{}

	if len(t.BottomButtons) > 0 {

		for _, b := range t.BottomButtons {
			if len(b.Texts) > 0 {
				// br := make([]tg.KeyboardButton, 0)
				br := tg.KeyboardButtonRow{}

				for _, t := range b.Texts {
					br.Buttons = append(br.Buttons, &tg.KeyboardButton{Text: t})
				}
				res.Rows = append(res.Rows, br)
			} else {
				res.Rows = append(res.Rows, tg.KeyboardButtonRow{
					Buttons: []tg.KeyboardButtonClass{
						&tg.KeyboardButton{Text: b.Text},
					},
				})
			}
		}

	}
	return &res
}

func (t *outcomingTextMessage) InlineKeyboardMarkup() *tg.ReplyInlineMarkup {
	// ReplyKeyboardRemove
	res := tg.ReplyInlineMarkup{}

	count := 0

	for _, row := range t.Buttons {
		// br := make([]tg.InlineKeyboardButton, 0)

		br := tg.KeyboardButtonRow{}

		for _, b := range row {
			count++
			br.Buttons = append(br.Buttons, &tg.KeyboardButtonCallback{
				Text: b.Text,
				Data: []byte(b.CallbackData()),
			})
			// br = append(br, tg.InlineKeyboardButton{
			// 	Text:         b.Text,
			// 	CallbackData: b.CallbackData(),
			// })
		}
		res.Rows = append(res.Rows, br)
	}

	if count == 0 {
		return nil
	}

	return &res
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

func equalReplyMarkup(a tg.ReplyMarkupClass, b tg.ReplyMarkupClass) bool {
	logger := DevLogger()

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
