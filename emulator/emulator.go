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

func FakeServerToInput(fakeServer *FakeServer) *DrawInput {
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

func DrawFakeServer(fakeServer *FakeServer) {
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

// func EmulateApplication[S any, A any](app tgbot.Application[S, A]) {
// 	// app.CreateAppState()
// 	// EmulatorDraw(&DrawInput{})
// }

/*

func main() {
	a := app.New()
	w := a.NewWindow("Emulator")

	inp := emul.DrawInput{
		Boxes: []emul.MessageBox{
			{
				Text: "Ты сдохнешь тупой гад!!!!",
				Buttons: emul.MessageButtons{
					Rows: []emul.ButtonsRow{
						{
							Butts: []emul.Butt{
								{Title: "Ударить", CallbackString: "Ударить"},
								{Title: "Убежать", CallbackString: "Убежать"},
							}},
					},
				},
			},
			{
				Text:    "Свин!!!",
				Buttons: emul.MessageButtons{},
			},
		},
	}

	output := emul.EmulatorDraw(
		inp, emul.ActionsHandler{
			CallbackHandlers: func(s string) {
				fmt.Printf("Button clicked: %s", s)
			},
			UserInputHandler: func(s string) {
				fmt.Printf("User input: %s", s)
			},
		},
	)

	wc := container.NewGridWrap(
		fyne.Size{Width: 300},
		container.NewStack(
			canvas.NewRectangle(color.Black),
			output,
		),
	)

	w.SetContent(container.NewCenter(wc))

	w.ShowAndRun()
}

*/
