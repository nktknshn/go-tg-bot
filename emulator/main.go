package emulator

import (
	"context"
	"image/color"
	"math/rand"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/go-telegram/bot/models"
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
				&models.Update{
					ID: int64(rand.Int()),
					CallbackQuery: &models.CallbackQuery{
						Data: s,
						Message: &models.Message{
							ID: rand.Int(),
							Chat: models.Chat{
								ID: chatID,
							},
						},
					},
				})

		},
		UserInputHandler: func(s string) {
			logger.Info("user input handler", zap.String("input", s))

			dispatcher.HandleUpdate(
				context.Background(),
				bot,
				&models.Update{
					ID: int64(rand.Int()),
					Message: &models.Message{
						ID:   rand.Int(),
						Text: s,
						Chat: models.Chat{
							ID: chatID,
						},
					},
				})
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
