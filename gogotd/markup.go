package gogotd

import (
	"encoding/json"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tg"
	"github.com/nktknshn/go-tg-bot/tgbot/logging"
	"go.uber.org/zap"
)

func EqualReplyMarkup(a tg.ReplyMarkupClass, b tg.ReplyMarkupClass) bool {
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

func EqualReplyKeyboardMarkup(a *tg.ReplyKeyboardMarkup, b *tg.ReplyKeyboardMarkup) bool {
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

func EqualInlineKeyboardButtonClass(a tg.KeyboardButtonClass, b tg.KeyboardButtonClass) bool {
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

func EqualInlineKeyboardMarkup(a *tg.ReplyInlineMarkup, b *tg.ReplyInlineMarkup) bool {

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

			if !EqualInlineKeyboardButtonClass(button, b.Rows[i].Buttons[j]) {
				return false
			}
		}
	}

	return true
}
