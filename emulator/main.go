package emulator

import (
	"context"
	"image/color"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	tgbot "github.com/nktknshn/go-tg-bot"
	"go.uber.org/zap"
)

func EmulatorMain(
	bot *FakeBot,
	dispatcher *tgbot.ChatsDispatcher,
) {
	a := app.New()
	w := a.NewWindow("Emulator")
	// bot := emulator.NewFakeBot()

	chatID := int64(1)

	handlers := ActionsHandler{
		CallbackHandlers: func(s string) {
			logger.Info("user callback handler", zap.String("input", s))

			dispatcher.HandleUpdate(
				context.Background(),
				bot,
				NewCallbackQueryUpdate(CallbackQueryUpdate{
					Data: s,
					UpdateProps: UpdateProps{
						ChatID: chatID,
					},
				}))

		},
		UserInputHandler: func(s string) {
			logger.Info("user input handler", zap.String("input", s))

			dispatcher.HandleUpdate(
				context.Background(),
				bot,
				NewTextMessageUpdate(TextMessageUpdate{
					Text: s,
					UpdateProps: UpdateProps{
						ChatID: chatID,
					},
				}))
		},
	}

	updateInterface := func() {

		output := EmulatorDraw(
			FakeServerToInput(bot),
			&handlers,
		)

		wc := container.NewGridWrap(
			fyne.Size{Width: 300},
			container.NewStack(
				canvas.NewRectangle(color.Black),
				output,
			),
		)

		w.SetContent(container.NewCenter(wc))
	}

	updateInterface()

	bot.SetUpdateCallback(func() {
		go updateInterface()
	})

	w.ShowAndRun()
}
