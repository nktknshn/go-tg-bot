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
	// Value1    int `tgbot:"ctx"`
	Context AppGlobalContext
}

func (a *TestAppFullCycle1Root) Render(o tgbot.OO) {
	o.Messagef("TestAppFullCycle1Root %v", a.Context.Value1)

	o.Comp(&TestAppFullCycle1Comp1{})
}

type TestAppFullCycle1Comp1 struct {
	Render1 bool `tgbot:"ctx"`
}

func (a *TestAppFullCycle1Comp1) Render(o tgbot.OO) {
	o.Message("TestAppFullCycle1Comp1")

	if a.Render1 {
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
	NestedValue bool `tgbot:"ctx"`
}

func (a *TestAppFullCycle1Comp1Comp2) Render(o tgbot.OO) {
	o.Message("TestAppFullCycle1Comp1Comp2")
}

var app = tgbot.NewApplication[TestAppFullCycle1State, any](
	// create initial state
	func(tc *tgbot.TelegramContext) TestAppFullCycle1State {
		return TestAppFullCycle1State{}
	},
	// create root component
	func(tafcs TestAppFullCycle1State) tgbot.Comp[any] {
		return &TestAppFullCycle1Root{}
	},
	func(ac *tgbot.ApplicationContext[TestAppFullCycle1State, any], tc *tgbot.TelegramContext, a any) {

	},
	&tgbot.NewApplicationProps[TestAppFullCycle1State, any]{
		CreateGlobalContext: func(ics *tgbot.InternalChatState[TestAppFullCycle1State, any]) tgbot.GlobalContext {
			ctx := tgbot.NewGlobalContext()

			ctx.Add("NestedValue", ics.AppState.NestedValue)
			ctx.Add("Value1", 1)
			ctx.Add("Render1", true)

			return ctx
		},
	},
)

func TestAppFullCycle1(t *testing.T) {
	t.Log("TestAppFullCycle")
	bot := emulator.NewFakeBot()

	app.NewHandler(&tgbot.TelegramContext{
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
