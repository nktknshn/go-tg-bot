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
	Counter int
	Error   error
}

type Action struct {
	Increment int
	Error     error
}

func App(props struct {
	counter int
	err     error
}) tgbot.Comp[Action] {
	return func(o tgbot.O[Action]) {

		o.InputHandler(func(s string) Action {

			if s == "/start" {
				return Action{}
			}

			v, err := strconv.Atoi(s)

			if err != nil {
				return Action{Error: err}
			}

			return Action{Increment: v}
		})

		o.Messagef("Counter value: %v", props.counter)
		if props.err != nil {
			o.Messagef("Error: %v", props.err)
		}

		o.Button("Increment", func() Action {
			return Action{Increment: 1}
		})
		o.Button("Decrement", func() Action {
			return Action{Increment: -1}
		})
	}
}

var counterApp = tgbot.NewApplication[State, Action](
	func(tc *tgbot.TelegramContext) State {
		tc.Logger.Info("CreateAppState")

		return State{Counter: 0}
	},
	func(ac *tgbot.ApplicationContext[State, Action], tc *tgbot.TelegramContext, a Action) {
		tc.Logger.Info("HandleAction", zap.Any("Increment", a.Increment))
		ac.State.AppState.Counter += a.Increment
	},
	func(s State) tgbot.Comp[Action] {
		return App(struct {
			counter int
			err     error
		}{
			counter: s.Counter,
			err:     nil,
		})
	},
	&tgbot.NewApplicationProps[State, Action]{},
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
