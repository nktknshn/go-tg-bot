package outcoming

import (
	"fmt"

	"github.com/gotd/td/tg"
	"github.com/nktknshn/go-tg-bot/gogotd"
	"github.com/nktknshn/go-tg-bot/tgbot/component"
)

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
		gogotd.EqualInlineKeyboardMarkup(t.InlineKeyboardMarkup(), oth.InlineKeyboardMarkup()) &&
		gogotd.EqualReplyKeyboardMarkup(t.ReplyKeyboardMarkup(), oth.ReplyKeyboardMarkup())
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
	res := tg.ReplyInlineMarkup{}

	count := 0

	for _, row := range t.Buttons {

		br := tg.KeyboardButtonRow{}

		for _, b := range row {
			count++
			br.Buttons = append(br.Buttons, &tg.KeyboardButtonCallback{
				Text: b.Text,
				Data: []byte(b.CallbackData()),
			})
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
