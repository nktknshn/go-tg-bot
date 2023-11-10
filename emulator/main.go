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

type UpdateProps struct {
	ChatID int64
	UserID int64
}

type CallbackQueryUpdate struct {
	Data string
	UpdateProps
}

func NewCallbackQueryUpdate(props CallbackQueryUpdate) *models.Update {
	return &models.Update{
		ID: int64(rand.Int()),
		CallbackQuery: &models.CallbackQuery{
			Data: props.Data,
			Message: &models.Message{
				ID: rand.Int(),
				Chat: models.Chat{
					ID: props.ChatID,
				},
				From: &models.User{
					ID:       int64(props.UserID),
					Username: "username",
				},
			},
		},
	}
}

type TextMessageUpdate struct {
	Text string
	UpdateProps
}

func NewTextMessageUpdate(props TextMessageUpdate) *models.Update {
	return &models.Update{
		ID: int64(rand.Int()),
		Message: &models.Message{
			ID:   rand.Int(),
			Text: props.Text,
			Chat: models.Chat{
				ID: props.ChatID,
			},
			From: &models.User{
				ID:       int64(props.UserID),
				Username: "username",
			},
		},
	}
}

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
