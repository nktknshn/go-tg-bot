package main

import (
	"context"
	"fmt"
	"image/color"

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

	},
	CreateAppState: func(tc *tgbot.TelegramContext) State {
		return State{
			Counter: 0,
		}
	},
	HandleAction: func(ac *tgbot.ApplicationContext[State, Action], tc *tgbot.TelegramContext, a Action) {
		ac.State.AppState.Counter += a.Increment
	},
	HandleMessage: func(ac *tgbot.ApplicationContext[State, Action], tc *tgbot.TelegramContext) {
		ac.State.InputHandler(tc)
	},
	HandleCallback: func(ac *tgbot.ApplicationContext[State, Action], tc *tgbot.TelegramContext) {
		act, err := ac.State.CallbackHandler(tc)

		if err != nil {
			logger.Error("Error in HandleCallback", zap.Error(err))
			return
		}

		ac.State.AppState.Counter += act.Increment
	},
	RenderFunc: func(ac *tgbot.ApplicationContext[State, Action]) []tgbot.RenderedElement {
		els := tgbot.ComponentToElements2(
			App(struct{ counter int }{
				counter: ac.State.AppState.Counter,
			}),
		)

		res, err := tgbot.ElementsToMessagesAndHandlers[Action](els)

		if err != nil {
			logger.Error("Error in RenderFunc", zap.Error(err))
		}

		if len(res.InputHandlers) == 0 {
			logger.Debug("No input handlers")
		} else {
			ac.State.InputHandler = res.InputHandlers[0].Handler
		}

		callbackHandlers := make(map[string]func() Action)

		for _, m := range res.OutcomingMessages {
			switch el := m.(type) {
			case *tgbot.OutcomingTextMessage[Action]:
				for _, row := range el.Buttons {
					for _, butt := range row {
						callbackHandlers[butt.Action] = butt.OnClick
					}
				}
			default:

			}
		}

		ac.State.CallbackHandler = func(tc *tgbot.TelegramContext) (Action, error) {
			key := tc.Update.CallbackQuery.Data

			if handler, ok := callbackHandlers[key]; ok {
				// ac.App.HandleAction(ac, tc, )
				return handler(), nil
			} else {
				err := fmt.Errorf("no handler for callback %v", key)
				logger.Error("No handler for callback", zap.String("key", key))
				return Action{}, err
			}
		}

		actions := tgbot.GetRenderActions(
			ac.State.RenderedElements,
			res.OutcomingMessages,
		)

		rendered, err := tgbot.RenderActions(context.Background(), ac.State.Renderer, actions)

		if err != nil {
			logger.Error("Error in RenderFunc", zap.Error(err))
			return []tgbot.RenderedElement{}
		}

		return rendered

		// return res.OutcomingMessages

	},
	CreateChatRenderer: func(tc *tgbot.TelegramContext) tgbot.ChatRenderer {
		// return emulator.NewFakeServer()
		return server
	},
}

func EmulatorMain() {
	a := app.New()
	w := a.NewWindow("Emulator")

	handlers := emulator.ActionsHandler{
		CallbackHandlers: func(s string) {
			logger.Info("callback handler", zap.String("callback", s))
		},
		UserInputHandler: func(s string) {
			logger.Info("user input handler", zap.String("input", s))
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
	ctx := context.Background()

	handler := counterApp.NewHandler(&tgbot.TelegramContext{
		ChatID: 0,
		Bot:    emulator.NewFakeBot(),
		Ctx:    context.Background(),
		Update: &models.Update{
			Message: &models.Message{
				ID:   0,
				Text: "/start",
			},
		},
		Logger: logger,
	})

	server.Message(ctx, &tgbot.ChatRendererMessageProps{
		Text: "Hello",
	})

	// emulator.DrawFakeServer(server)
	EmulatorMain()
}
