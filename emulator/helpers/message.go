package helpers

import (
	"encoding/json"

	"github.com/gotd/td/tg"
)

type ReplyMarkupInlineSimple = [][]ButtonSimpl

type MessageSimple struct {
	Message string
	Buttons ReplyMarkupInlineSimple
}

func (mr MessageSimple) ToJson() string {
	json, err := json.Marshal(mr)

	if err != nil {
		panic(err)
	}

	return string(json)
}

type ButtonSimpl struct {
	Text string
	Data string
}

func ReplyMarkupAsSimple(rm tg.ReplyMarkupClass) ReplyMarkupInlineSimple {

	m, ok := rm.(*tg.ReplyInlineMarkup)

	if !ok {
		return nil
	}

	if m == nil {
		return nil
	}

	var buttons [][]ButtonSimpl

	for _, row := range m.Rows {
		var buttonsRow []ButtonSimpl
		for _, button := range row.Buttons {
			button, ok := button.(*tg.KeyboardButtonCallback)

			if !ok {
				continue
			}

			buttonsRow = append(buttonsRow, ButtonSimpl{
				Text: button.Text,
				Data: string(button.Data),
			})
		}

		buttons = append(buttons, buttonsRow)
	}

	return buttons
}

func MessageAsSimple(msg *tg.Message) MessageSimple {
	return MessageSimple{
		Message: msg.Message,
		Buttons: ReplyMarkupAsSimple(msg.ReplyMarkup),
	}
}

func MessageAsJson(msg *tg.Message) string {
	raw := MessageAsSimple(msg)

	return raw.ToJson()
}
