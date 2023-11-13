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
		o.Message("Night")
	} else {
		tgbot.GetLogger().Debug("day")
		o.Message("Day")
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

func TestRunCreateElements(t *testing.T) {
	comp := App1{Counter: 1}

	res := tgbot.CreateElements[any](&comp, nil)

	if len(res.Elements) != 4 {
		t.Fatal("len(res.Elements) != 4")
	}

	t.Logf("res: %s", res)
	t.Logf("Local Value: %v", res.TreeState.LocalStateTree.LocalStateClosure)

}

func TestLocalStateTree(t *testing.T) {
	type V struct {
		v int
	}
	ls := tgbot.LocalStateTree{
		LocalStateClosure: &tgbot.LocalStateClosure[any]{
			Initialized: true,
			Value:       1,
		},
		Children: &[]*tgbot.LocalStateTree{
			nil,
			{
				LocalStateClosure: &tgbot.LocalStateClosure[any]{
					Initialized: true,
					Value:       2,
				},
				Children: &[]*tgbot.LocalStateTree{},
			},
			{
				LocalStateClosure: &tgbot.LocalStateClosure[any]{
					Initialized: true,
					Value:       3,
				},
				Children: &[]*tgbot.LocalStateTree{},
			},
			{
				LocalStateClosure: &tgbot.LocalStateClosure[any]{
					Initialized: true,
					Value:       4,
				},
				Children: &[]*tgbot.LocalStateTree{
					{
						LocalStateClosure: &tgbot.LocalStateClosure[any]{
							Initialized: true,
							Value:       V{v: 5},
						},
						Children: &[]*tgbot.LocalStateTree{},
					},
				},
			},
		}}

	if ls.Get([]int{}).Value != 1 {
		t.Errorf("ls.Get([]int{}).Value != 1")
	}

	if ls.Get([]int{1}).Value != 2 {
		t.Errorf("ls.Get([]int{1}).Value != 2")
	}

	if ls.Get([]int{3}).Value != 4 {
		t.Errorf("ls.Get([]int{3}).Value != 4")
	}

	if ls.Get([]int{3, 0}).Value.(V).v != 5 {
		t.Errorf("ls.Get([]int{3, 0}).Value != 5")
	}

	ls.Set([]int{3, 0}, func(v any) any {
		return V{v: v.(V).v + 1}
	})

	if ls.Get([]int{3, 0}).Value.(V).v != 6 {
		t.Errorf("ls.Get([]int{3, 0}).Value != 6")
	}

	t.Logf("ls: %v", ls)

}
