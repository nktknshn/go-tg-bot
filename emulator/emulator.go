package emulator

import (
	"slices"

	"github.com/go-telegram/bot/models"
	tgbot "github.com/nktknshn/go-tg-bot"
	"golang.org/x/exp/maps"
)

func InlineKeyboardToButtons(ik *models.InlineKeyboardMarkup) MessageButtons {
	mb := MessageButtons{}

	for _, row := range ik.InlineKeyboard {
		mbrow := ButtonsRow{}

		for _, butt := range row {
			mbrow.Butts = append(mbrow.Butts, Butt{
				Title:          butt.Text,
				CallbackString: butt.CallbackData,
			})
		}

		mb.Rows = append(mb.Rows, mbrow)
	}

	return mb
}

func FakeServerToInput(fakeServer *FakeBot) *DrawInput {
	result := &DrawInput{}
	messages := maps.Values(fakeServer.Messages)

	slices.SortFunc(messages, func(a, b *models.Message) int {
		return a.ID - b.ID
	})

	for _, m := range fakeServer.Messages {

		mbs := InlineKeyboardToButtons(&m.ReplyMarkup)

		result.Boxes = append(result.Boxes, MessageBox{
			Text:    m.Text,
			Buttons: mbs,
		})
	}

	return result
}

var logger = tgbot.GetLogger()
