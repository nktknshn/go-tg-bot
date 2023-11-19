package tgbot_test

import (
	"fmt"
	"testing"

	tgbot "github.com/nktknshn/go-tg-bot"
)

type TestRunComponent1State struct {
	TrilpleSix int
}

type TestRunComponent1Comp struct {
	State tgbot.State[TestRunComponent1State]
}

func (c TestRunComponent1Comp) Render(o tgbot.OO) {
	state := c.State.Init(TestRunComponent1State{
		TrilpleSix: 666,
	})

	o.Messagef("Hello, world! %v", state.Get().TrilpleSix)
}

func TestRunComponent1(t *testing.T) {
	comp := TestRunComponent1Comp{}

	localStateClosure := tgbot.LocalStateClosure[any]{
		Initialized: true,
		Value: TestRunComponent1State{
			111,
		},
	}
	state := tgbot.State[any]{
		LocalStateClosure: &localStateClosure,
		Index:             []int{},
	}

	els, closure, ctx := tgbot.RunComponent(
		tgbot.GetLogger(), comp, tgbot.NewGlobalContextTyped(1), state,
	)

	fmt.Println("els", els)
	fmt.Println("closure", closure)
	fmt.Println("ctx", ctx)
}
