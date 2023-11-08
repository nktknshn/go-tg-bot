package main

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
	"github.com/nktknshn/go-tg-bot/emulator"
	"go.uber.org/zap"
)

type State struct {
	Counter int
}

type Action struct {
	Increment int
}

func App(props struct{ counter int }) tgbot.Comp[Action] {
	return func(o tgbot.O[Action]) {

		o.Messagef("Counter value: %v", props.counter)
		o.Button("Increment", func() Action {
			return Action{Increment: 1}
		})
		o.Button("Decrement", func() Action {
			return Action{Increment: -1}
		})
	}
}

var logger = tgbot.GetLogger()

var server = emulator.NewFakeServer()

var counterApp = &tgbot.Application[State, Action]{
	HandleInit: func(tc *tgbot.TelegramContext) {
		tc.Logger.Info("Init")
	},
	CreateAppState: func(tc *tgbot.TelegramContext) State {
		tc.Logger.Info("CreateAppState")

		return State{
			Counter: 0,
		}
	},
	HandleAction: func(ac *tgbot.ApplicationContext[State, Action], tc *tgbot.TelegramContext, a Action) {
		tc.Logger.Info("HandleAction", zap.Any("Increment", a.Increment))

		ac.State.AppState.Counter += a.Increment
	},
	HandleMessage: func(ac *tgbot.ApplicationContext[State, Action], tc *tgbot.TelegramContext) {
		tc.Logger.Info("HandleMessage", zap.Any("text", tc.Update.Message.Text))

		ac.State.InputHandler(tc)
	},
	HandleCallback: func(ac *tgbot.ApplicationContext[State, Action], tc *tgbot.TelegramContext) {
		tc.Logger.Info("HandleCallback", zap.Any("data", tc.Update.CallbackQuery.Data))

		act, err := ac.State.CallbackHandler(tc)

		if err != nil {
			logger.Error("Error in HandleCallback", zap.Error(err))
			return
		}

		ac.State.AppState.Counter += act.Increment
	},
	StateToComp: func(s State) tgbot.Comp[Action] {
		return App(struct{ counter int }{
			counter: s.Counter,
		})
	},
	RenderFunc: func(ac *tgbot.ApplicationContext[State, Action]) error {
		logger.Info("RenderFunc")

		res, err := ac.App.PreRender(ac)

		if err != nil {
			logger.Error("Error in RenderFunc", zap.Error(err))
			return err
		}

		rendered, err := res.ExecuteRender(ac.State.Renderer)

		if err != nil {
			logger.Error("Error in RenderFunc", zap.Error(err))
			return err
		}

		ac.State = &res.InternalChatState
		ac.State.RenderedElements = rendered

		return nil
	},
	CreateChatRenderer: func(tc *tgbot.TelegramContext) tgbot.ChatRenderer {
		// return emulator.NewFakeServer()
		return server
	},
}

func EmulatorMain(
	dispatcher *tgbot.ChatsDispatcher,
) {
	a := app.New()
	w := a.NewWindow("Emulator")
	bot := emulator.NewFakeBot()
	chatID := int64(1)

	handlers := emulator.ActionsHandler{
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

	output := emulator.EmulatorDraw(
		emulator.FakeServerToInput(server),
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
	w.ShowAndRun()
}

func main() {
	// ctx := context.Background()

	dispatcher := tgbot.NewChatsDispatcher(&tgbot.ChatsDispatcherProps{
		ChatFactory: func(tc *tgbot.TelegramContext) tgbot.ChatHandler {
			return counterApp.NewHandler(tc)
		},
	})

	EmulatorMain(dispatcher)

}
