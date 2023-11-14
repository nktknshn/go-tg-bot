package tgbot_test

import (
	"testing"

	tgbot "github.com/nktknshn/go-tg-bot"
)

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
