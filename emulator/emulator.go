package emulator

import (
	"github.com/go-telegram/bot/models"
	tgbot "github.com/nktknshn/go-tg-bot"
	"go.uber.org/zap"
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

func DrawFakeServer(fakeServer *FakeBot) {
	drawInput := FakeServerToInput(fakeServer)

	EmulatorDraw(drawInput, &ActionsHandler{
		CallbackHandlers: func(s string) {
			logger.Info("callback handler", zap.String("callback", s))
		},
		UserInputHandler: func(s string) {
			logger.Info("user input handler", zap.String("input", s))
		},
	})
}
