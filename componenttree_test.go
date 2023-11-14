package tgbot_test

import (
	"testing"

	tgbot "github.com/nktknshn/go-tg-bot"
)

type App1State struct {
	night bool
	hour  int
}

type App1 struct {
	// props
	Counter int
	// local state
	State tgbot.GetSetLocalStateImpl[App1State]
}

func (a *App1) Render(o tgbot.OO) {
	// TODO make it INIT
	lsgs := a.State.Init(App1State{
		night: false,
		hour:  3,
	})

	o.Message("Hello")

	if lsgs.Get().night {
		tgbot.GetLogger().Debug("night")
		o.Messagef("Night: %v", lsgs.Get().hour)
	} else {
		tgbot.GetLogger().Debug("day")
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

func TestRunComponent(t *testing.T) {
	comp := App1{Counter: 1}

	tgbot.RunComponent(&comp, tgbot.GetSetLocalStateImpl[any]{
		LocalState: tgbot.LocalStateClosure[any]{
			Initialized: true,
			Value:       App1State{night: true, hour: 2},
		},
		Index: []int{0},
	})

}

func TestRunCreateElements1(t *testing.T) {
	comp := App1{Counter: 1}

	res := tgbot.CreateElements[any](&comp, nil)

	if len(res.Elements) != 4 {
		t.Fatal("len(res.Elements) != 4")
	}

	if res.Elements[1].(*tgbot.ElementMessage).Text != "Day: 3" {
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

	res = tgbot.CreateElements[any](&comp, &res.TreeState)

	if len(res.Elements) != 4 {
		t.Fatal("len(res.Elements) != 4")
	}

	if res.Elements[1].(*tgbot.ElementMessage).Text != "Night: 8" {
		t.Fatal("Night: 8 was expected")
	}

}

type TestNestedCompApp struct{}

func (c *TestNestedCompApp) Render(o tgbot.OO) {
	o.Message("App1")

	o.Comp(&TestNestedComp1{})
}

type Context[S any] interface {
	Get() S
}

type TestNestedComp1Context struct {
	Flag1 bool
}

type TestNestedComp1 struct {
	// props
	Flag1 bool
	// requests data from the global state
	// ctx      Context[TestNestedComp1Context]
	CtxFlag1 bool
}

func (c *TestNestedComp1) Render(o tgbot.OO) {
	o.Message("Comp1")

	if c.Flag1 {
		o.Comp(&TestNestedComp2{})
	} else {
		o.Comp(&TestNestedComp3{})
	}
}

type TestNestedComp2 struct{}

func (c *TestNestedComp2) Render(o tgbot.OO) {
	o.Message("Comp2")

}

type TestNestedComp3 struct{}

func (c *TestNestedComp3) Render(o tgbot.OO) {
	o.Message("Comp3")
}

func TestNestedComp(t *testing.T) {
	res := tgbot.CreateElements(&TestNestedCompApp{}, nil)
	res = tgbot.CreateElements(&TestNestedCompApp{}, &res.TreeState)

	t.Log(res.Elements)
}
