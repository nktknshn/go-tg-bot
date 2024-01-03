package outcoming

import (
	"encoding/json"
	"fmt"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tg"
	"github.com/nktknshn/go-tg-bot/tgbot/component"
	"github.com/nktknshn/go-tg-bot/tgbot/logging"
	"go.uber.org/zap"
)

// Messages that are going to be sent to user
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
	ElementFile component.ElementFile
	Message     *tg.Message
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
	ElementUserMessage component.ElementUserMessage
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
	ElementPhotoGroup component.ElementPhotoGroup
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

type OutcomingTextMessage struct {
	Text          string
	Buttons       [][]component.ElementButton
	BottomButtons []component.ElementBottomButton
	isComplete    bool
	// TODO RequestLocation
}

func (t *OutcomingTextMessage) String() string {
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

func NewOutcomingTextMessage(text string) *OutcomingTextMessage {
	buttons := make([][]component.ElementButton, 0)
	buttons = append(buttons, make([]component.ElementButton, 0))

	return &OutcomingTextMessage{
		Text:       text,
		Buttons:    buttons,
		isComplete: false,
	}
}

func (t *OutcomingTextMessage) OutcomingKind() string {
	return KindOutcomingTextMessage
}

func (t *OutcomingTextMessage) Equal(other OutcomingMessage) bool {
	if other.OutcomingKind() != KindOutcomingTextMessage {
		return false
	}

	oth := other.(*OutcomingTextMessage)

	return t.Text == oth.Text &&
		equalInlineKeyboardMarkup(t.InlineKeyboardMarkup(), oth.InlineKeyboardMarkup()) &&
		equalReplyKeyboardMarkup(t.ReplyKeyboardMarkup(), oth.ReplyKeyboardMarkup())
}

func (t *OutcomingTextMessage) ConcatText(text string) {
	t.Text += "\n" + text
}

func (t *OutcomingTextMessage) SetComplete() {
	t.isComplete = true
}

func (t *OutcomingTextMessage) ReplyMarkup() tg.ReplyMarkupClass {

	if len(t.BottomButtons) > 0 {
		return t.ReplyKeyboardMarkup()
	}

	return t.InlineKeyboardMarkup()
}

func (t *OutcomingTextMessage) ReplyKeyboardMarkup() *tg.ReplyKeyboardMarkup {
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

func (t *OutcomingTextMessage) InlineKeyboardMarkup() *tg.ReplyInlineMarkup {
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

func (t *OutcomingTextMessage) AddButton(button *component.ElementButton) {
	emptyButtons := len(t.Buttons) == 1 && len(t.Buttons[0]) == 0

	if button.NextRow && !emptyButtons {
		t.Buttons = append(t.Buttons, make([]component.ElementButton, 0))
	}

	t.Buttons[len(t.Buttons)-1] = append(t.Buttons[len(t.Buttons)-1], *button)
}

func (t *OutcomingTextMessage) AddButtonsRow(buttonsRow *component.ElementButtonsRow) {
	t.Buttons = append(t.Buttons, buttonsRow.Buttons())
}

func equalReplyMarkup(a tg.ReplyMarkupClass, b tg.ReplyMarkupClass) bool {
	logger := logging.Logger()

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
