package examples

import (
	"fmt"

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

var app = &tgbot.Application[State, Action]{
	HandleInit: func(tc *tgbot.TelegramContext) {

	},
	CreateAppState: func(tc *tgbot.TelegramContext) State {
		return State{}
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

		var callbackHandlers map[string]func() Action

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

		return res.OutcomingMessages

	},
	CreateChatRenderer: func(tc *tgbot.TelegramContext) tgbot.ChatRenderer {
		return emulator.NewFakeServer()
	},
}

func main() {

}
