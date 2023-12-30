package tgbot

import (
	"testing"
)

type App1State struct {
	night bool
	hour  int
}

type App1 struct {
	// props
	Counter int
	// local state
	State CompState[App1State]
}

func (a *App1) Render(o O) {
	// TODO make it INIT
	lsgs := a.State.Init(App1State{
		night: false,
		hour:  3,
	})

	o.Message("Hello")

	if lsgs.Get().night {
		GetLogger().Debug("night")
		o.Messagef("Night: %v", lsgs.Get().hour)
	} else {
		GetLogger().Debug("day")
		o.Messagef("Day: %v", lsgs.Get().hour)
	}

	o.Button("Toggle Day/Ngiht", func() any {
		return lsgs.Set(func(as App1State) App1State {
			// as.boolean = !as.boolean
			return App1State{night: !as.night}
		})
	})

	o.Messagef("Counter: %v", a.Counter)
}

type EmptyContext struct{}

func TestRunComponent(t *testing.T) {
	comp := App1{Counter: 1}
	globalContext := newGlobalContextTyped[any](EmptyContext{})

	runComponent(
		GetLogger(),
		&comp,
		globalContext,
		CompState[any]{
			LocalStateClosure: &localStateClosure[any]{
				Initialized: true,
				Value:       App1State{night: true, hour: 2},
			},
			Index: []int{0},
		})

}

func TestRunCreateElements1(t *testing.T) {

	globalContext := newGlobalContextTyped[any](EmptyContext{})

	comp := App1{Counter: 1}

	res := createElements(&comp, globalContext, nil)

	if len(res.Elements) != 4 {
		t.Fatal("len(res.Elements) != 4")
	}

	if res.Elements[1].(*elementMessage).Text != "Day: 3" {
		t.Fatal("Day: 3 was expected")
	}

	// t.Logf("res: %s", res)
	// t.Logf("Local Value: %v", res.TreeState.LocalStateTree.LocalStateClosure)

	res.TreeState.LocalStateTree.Set([]int{}, func(a any) any {
		return App1State{
			night: !a.(App1State).night,
			hour:  a.(App1State).hour + 5,
		}
	})

	res = createElements(&comp, globalContext, &res.TreeState)

	if len(res.Elements) != 4 {
		t.Fatal("len(res.Elements) != 4")
	}

	if res.Elements[1].(*elementMessage).Text != "Night: 8" {
		t.Fatal("Night: 8 was expected")
	}

}

type TestNestedCompContext struct {
	Flag1 bool
}

type TestNestedCompApp struct{}

func (c *TestNestedCompApp) Render(o O) {
	o.Message("App1")

	o.Comp(&TestNestedComp1{})
}

type Context[S any] interface {
	Get() S
}

type TestNestedComp1 struct {
	Context TestNestedCompContext
}

func (c *TestNestedComp1) Render(o O) {
	o.Message("Comp1")

	if c.Context.Flag1 {
		o.Comp(&TestNestedComp2{})
	} else {
		o.Comp(&TestNestedComp3{})
	}
}

type TestNestedComp2 struct{}

func (c *TestNestedComp2) Render(o O) {
	o.Message("Comp2")

}

type TestNestedComp3 struct{}

func (c *TestNestedComp3) Render(o O) {
	o.Message("Comp3")
}

func TestNestedComp(t *testing.T) {

	globalContext := newGlobalContextTyped[any](TestNestedCompContext{
		Flag1: false,
	})

	t.Log("globalContext", globalContext)

	res := createElements(&TestNestedCompApp{}, globalContext, nil)

	globalContext = newGlobalContextTyped[any](TestNestedCompContext{
		Flag1: true,
	})

	res = createElements(&TestNestedCompApp{}, globalContext, &res.TreeState)

	t.Log(res.Elements)
}
