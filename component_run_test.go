package tgbot

import (
	"fmt"
	"testing"
)

type TestRunComponent1State struct {
	TrilpleSix int
}

type TestRunComponent1Comp struct {
	State CompState[TestRunComponent1State]
}

func (c TestRunComponent1Comp) Render(o O) {
	state := c.State.Init(TestRunComponent1State{
		TrilpleSix: 666,
	})

	o.Messagef("Hello, world! %v", state.Get().TrilpleSix)
}

func TestRunComponent1(t *testing.T) {
	comp := TestRunComponent1Comp{}

	localStateClosure := localStateClosure[any]{
		Initialized: true,
		Value: TestRunComponent1State{
			111,
		},
	}
	state := CompState[any]{
		LocalStateClosure: &localStateClosure,
		Index:             []int{},
	}

	els, closure, ctx := runComponent(
		GetLogger(), comp, newGlobalContextTyped(1), state,
	)

	fmt.Println("els", els)
	fmt.Println("closure", closure)
	fmt.Println("ctx", ctx)
}
