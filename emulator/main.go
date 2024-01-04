package emulator

import (
	"image/color"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/nktknshn/go-tg-bot/tgbot/dispatcher"
	"go.uber.org/zap"
)

func EmulatorMain(
	bot *FakeBot,
	dispatcher *dispatcher.ChatsDispatcher,
) {
	a := app.New()
	w := a.NewWindow("Emulator")

	emul := NewEmulator()

	// chatID := int64(1)

	handlers := ActionsHandler{
		CallbackHandlers: func(s string) {
			logger.Info("user callback handler", zap.String("input", s))

			// dispatcher.HandleUpdate(
			// 	context.Background(),
			// 	bot,
			// 	NewCallbackQueryUpdate(CallbackQueryUpdate{
			// 		Data: s,
			// 		UpdateProps: UpdateProps{
			// 			ChatID: chatID,
			// 		},
			// 	}))

		},
		UserInputHandler: func(s string) {
			logger.Info("user input handler", zap.String("input", s))
			// bot.AddUserMessage

			// dispatcher.HandleUpdate(
			// 	context.Background(),
			// 	bot,
			// 	userMessageUpdate)
		},
	}

	emul.SetHandler(&handlers)

	updateInterface := func() {

		output := emul.Draw(
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

	// bot.SetReplyCallback(func() {
	// 	emul.SetCallbackReceived()
	// })

	emul.SetEmulatorStateUpdatedCallback(func() {
		logger.Info("emulator state updated")
		go updateInterface()
	})

	w.ShowAndRun()
}
