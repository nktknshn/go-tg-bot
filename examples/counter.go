package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	tgbot "github.com/nktknshn/go-tg-bot"
	"github.com/nktknshn/go-tg-bot/emulator"
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

func AppInputHandler(o tgbot.OO) {

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
	Username string
	state    tgbot.LocalStateIniter[WelcomState]
}

func (w *Welcom) Render(o tgbot.OO) {

	// ls := w.state.Init(WelcomState{})

	// if ls.Get().hideName {
	// 	o.Message("Welcome")
	// } else {
	o.MessagePartf("Welcome %v", w.Username)
	o.MessagePart("/hide_name to hide your name")
	o.MessageComplete()
	// }

	// o.InputHandler(func(s string) any {
	// 	if s == "/hide_name" {
	// 		return ls.Set(func(s WelcomState) {
	// 			s.hideName = true
	// 		})
	// 	}

	// 	return nil
	// })
}

type App struct {
	Props
}

func (app *App) Render(o tgbot.OO) {

	AppInputHandler(o)

	o.Comp(&Welcom{
		Username: app.Username,
	})

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

var counterApp = tgbot.NewApplication[State, any](
	func(tc *tgbot.TelegramContext) State {
		// tc.Logger.Info("CreateAppState")
		uname := fmt.Sprintf("%v (%v)", tc.Update.Message.From.Username, tc.Update.Message.From.ID)

		return State{Counter: 0, Username: uname}
	},
	func(ac *tgbot.ApplicationContext[State, any], tc *tgbot.TelegramContext, a any) {
		// tc.Logger.Info("HandleAction", zap.Any("action", a))

		switch a := a.(type) {
		case ActionReload:
			ac.State.RenderedElements = make([]tgbot.RenderedElement, 0)
		case ActionCounter:
			ac.State.AppState.Counter += a.Increment
			ac.State.AppState.Error = a.Error
		}

	},
	func(s State) tgbot.Comp[any] {
		app := App{Props(s)}

		return &app
	},
)

func runEmulator(dispatcher *tgbot.ChatsDispatcher) {
	bot := emulator.NewFakeBot()
	emulator.EmulatorMain(bot, dispatcher)
}

func runReal(dispatcher *tgbot.ChatsDispatcher) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	token := os.Getenv("BOT_TOKEN")

	if token == "" {
		logger.Fatal("BOT_TOKEN env variable is not set")
		os.Exit(1)
	}

	bot, err := bot.New(token, bot.WithDefaultHandler(func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		dispatcher.HandleUpdate(ctx, bot, update)
	}))

	if err != nil {
		logger.Fatal("Error creating bot", zap.Error(err))
		os.Exit(1)
	}

	bot.Start(ctx)

}

var logger = tgbot.GetLogger()

func main() {
	// if first argument is "emulator", run emulator
	flag.Parse()

	logger.Debug("Starting bot", zap.Any("args", flag.Args()))

	dispatcher := tgbot.NewChatsDispatcher(&tgbot.ChatsDispatcherProps{
		ChatFactory: func(tc *tgbot.TelegramContext) tgbot.ChatHandler {
			return counterApp.NewHandler(tc)
		},
	})

	if len(flag.Args()) > 0 && flag.Args()[0] == "emulator" {
		runEmulator(dispatcher)
	} else if len(flag.Args()) > 0 && flag.Args()[0] == "real" {
		runReal(dispatcher)
	} else {
		logger.Fatal("Unknown argument", zap.Any("args", flag.Args()))
		fmt.Println("emulator or real")
		os.Exit(1)
	}
}
