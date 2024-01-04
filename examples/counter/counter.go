package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/nktknshn/go-tg-bot/emulator"
	tgbot "github.com/nktknshn/go-tg-bot/tgbot"
	"github.com/nktknshn/go-tg-bot/tgbot/application"
	"github.com/nktknshn/go-tg-bot/tgbot/component"
	"github.com/nktknshn/go-tg-bot/tgbot/dispatcher"
	"github.com/nktknshn/go-tg-bot/tgbot/gotd"
	"github.com/nktknshn/go-tg-bot/tgbot/logging"
	"go.uber.org/zap"
)

type State struct {
	Counter  int
	Error    error
	Username string
}

type ActionReload struct{}

type ActionCounter struct {
	Increment int
	Error     error
}

type Props struct {
	Counter  int
	Error    error
	Username string
}

func AppInputHandler(o tgbot.O) {

	o.InputHandler(func(s string) any {
		if s == "/start" {
			return &ActionReload{}
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			return &ActionCounter{Error: err}
		}
		return &ActionCounter{Increment: v}
	})

}

type WelcomState struct {
	hideName bool
}

type Welcom struct {
	Username string `tgbot:"ctx"`
	State    component.State[WelcomState]
}

func (w *Welcom) Render(o tgbot.O) {

	ls := w.State.Init(WelcomState{})

	if ls.Get().hideName {
		o.Message("Welcome")
	} else {
		o.MessagePartf("Welcome %v", w.Username)
		o.MessagePart("/hide_name to hide your name")
		o.MessageComplete()
	}

	o.InputHandler(func(s string) any {
		if s == "/hide_name" {
			return ls.Set(func(s WelcomState) WelcomState {
				return WelcomState{
					hideName: s.hideName,
				}
			})
		}

		return nil
	})
}

type App struct {
	Props
}

func (app *App) Render(o tgbot.O) {

	AppInputHandler(o)

	o.Comp(&Welcom{})

	o.Messagef("Counter value: %v", app.Counter)

	if app.Error != nil {
		o.Messagef("Error: %v", app.Error)
	}

	o.Button("Increment", func() any {
		return ActionCounter{Increment: 1}
	})
	o.Button("Decrement", func() any {
		return ActionCounter{Increment: -1}
	})
}

var counterApp = application.New[State, any](
	func(tc *tgbot.TelegramUpdateContext) State {
		// tc.Logger.Info("CreateAppState")
		// tc.Message.
		// uname := fmt.Sprintf("%v (%v)", tc.Message.From.Username, tc.Update.Message.From.ID)
		var uname string = "unknown"

		if tc, ok := tc.AsTextMessage(); ok {
			uname = fmt.Sprintf("%v", tc.Message.PeerID.String())
		} else if tc, ok := tc.AsCallback(); ok {
			uname = fmt.Sprintf("%v", tc.UpdateBotCallbackQuery.Peer.String())
		}

		return State{Counter: 0, Username: uname}
	},
	func(s State) tgbot.Comp {
		app := App{Props(s)}

		return &app
	},
	func(ac *application.ApplicationChat[State, any], tc *tgbot.TelegramUpdateContext, a any) {
		// tc.Logger.Info("HandleAction", zap.Any("action", a))

		switch a := a.(type) {
		case ActionReload:
			ac.State.ResetRenderedElements()
		case ActionCounter:
			ac.State.AppState.Counter += a.Increment
			ac.State.AppState.Error = a.Error
		}
	},
)

var logger = logging.DevLogger()

func main() {
	// if first argument is "emulator", run emulator
	flag.Parse()

	ctx := context.Background()
	logger.Debug("Starting bot", zap.Any("args", flag.Args()))

	dispatcher := dispatcher.ForApplication(counterApp)

	if len(flag.Args()) > 0 && flag.Args()[0] == "emulator" {
		emulator.Run(logger, dispatcher)
	} else if len(flag.Args()) > 0 && flag.Args()[0] == "real" {
		gotd.Run(ctx, logger, dispatcher)
	} else {
		logger.Fatal("Unknown argument", zap.Any("args", flag.Args()))
		fmt.Println("emulator or real")
		os.Exit(1)
	}
}
