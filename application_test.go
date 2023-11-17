package tgbot_test

import (
	"context"
	"testing"

	tgbot "github.com/nktknshn/go-tg-bot"
	"github.com/nktknshn/go-tg-bot/emulator"
)

type AppGlobalContext struct {
	Value1      int
	Render1     bool
	NestedValue bool
}

/*
Here we are going to test how all the components work together.
*/

type TestAppFullCycle1State struct {
	StateA      string
	StateB      int
	StateC      []int
	NestedValue bool
}

type TestAppFullCycle1Root struct {
	RootProps string
	Context   AppGlobalContext
}

func (a *TestAppFullCycle1Root) Render(o tgbot.OO) {
	o.Messagef("TestAppFullCycle1Root %v", a.Context.Value1)

	o.Comp(&TestAppFullCycle1Comp1{})
}

type TestAppFullCycle1Comp1 struct {
	Context AppGlobalContext
}

func (a *TestAppFullCycle1Comp1) Render(o tgbot.OO) {
	o.Message("TestAppFullCycle1Comp1")

	if a.Context.Render1 {
		o.Comp(&TestAppFullCycle1Comp1Comp1{})
	} else {
		o.Comp(&TestAppFullCycle1Comp1Comp2{})
	}
}

type TestAppFullCycle1Comp1Comp1 struct{}

func (a *TestAppFullCycle1Comp1Comp1) Render(o tgbot.OO) {
	o.Message("TestAppFullCycle1Comp1Comp1")
}

type TestAppFullCycle1Comp1Comp2 struct {
	NestedValue bool
}

func (a *TestAppFullCycle1Comp1Comp2) Render(o tgbot.OO) {
	o.Message("TestAppFullCycle1Comp1Comp2")
}

var app = tgbot.NewApplication[TestAppFullCycle1State, any, AppGlobalContext](
	// create initial state
	func(tc *tgbot.TelegramContext) TestAppFullCycle1State {
		return TestAppFullCycle1State{}
	},
	// create root component
	func(tafcs TestAppFullCycle1State) tgbot.Comp[any] {
		return &TestAppFullCycle1Root{}
	},
	// handle action
	func(ac *tgbot.ApplicationContext[TestAppFullCycle1State, any, AppGlobalContext], tc *tgbot.TelegramContext, a any) {

	},
	// create global context
	&tgbot.NewApplicationProps[TestAppFullCycle1State, any, AppGlobalContext]{
		CreateGlobalContext: func(ics *tgbot.InternalChatState[TestAppFullCycle1State, any, AppGlobalContext]) tgbot.GlobalContextTyped[AppGlobalContext] {

			ctx := tgbot.NewGlobalContextTyped(AppGlobalContext{
				Value1:      1,
				Render1:     true,
				NestedValue: false,
			})

			return ctx
		},
	},
)

func TestAppFullCycle1(t *testing.T) {
	t.Log("TestAppFullCycle")
	bot := emulator.NewFakeBot()

	handler := app.NewHandler(&tgbot.TelegramContext{
		ChatID: 1,
		Bot:    bot,
		Ctx:    context.Background(),
		Update: emulator.NewTextMessageUpdate(emulator.TextMessageUpdate{
			Text: "test",
			UpdateProps: emulator.UpdateProps{
				ChatID: 1,
				UserID: 1,
			},
		}),
		Logger: tgbot.GetLogger(),
	})

	if len(bot.Messages) != 0 {
		t.Fatal("Expected empty")
	}

	handler.HandleUpdate(&tgbot.TelegramContext{
		ChatID: 1,
		Bot:    bot,
		Ctx:    context.Background(),
		Update: emulator.NewTextMessageUpdate(emulator.TextMessageUpdate{
			Text: "test",
			UpdateProps: emulator.UpdateProps{
				ChatID: 1,
				UserID: 1,
			},
		}),
		Logger: tgbot.GetLogger(),
	})
}
