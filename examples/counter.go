package main

import (
	"flag"
	"fmt"
	"strconv"

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

var logger = tgbot.GetLogger()

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

func main() {
	fmt.Println(flag.Arg(0))

	dispatcher := tgbot.NewChatsDispatcher(&tgbot.ChatsDispatcherProps{
		ChatFactory: func(tc *tgbot.TelegramContext) tgbot.ChatHandler {
			return counterApp.NewHandler(tc)
		},
	})

	bot := emulator.NewFakeBot()

	emulator.EmulatorMain(bot, dispatcher)

}
