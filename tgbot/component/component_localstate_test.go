package component

import (
	"testing"
)

func TestLocalStateTree(t *testing.T) {
	type V struct {
		v int
	}
	ls := localStateTree{
		LocalStateClosure: &localStateClosure[any]{
			Initialized: true,
			Value:       1,
		},
		Children: &[]*localStateTree{
			nil,
			{
				LocalStateClosure: &localStateClosure[any]{
					Initialized: true,
					Value:       2,
				},
				Children: &[]*localStateTree{},
			},
			{
				LocalStateClosure: &localStateClosure[any]{
					Initialized: true,
					Value:       3,
				},
				Children: &[]*localStateTree{},
			},
			{
				LocalStateClosure: &localStateClosure[any]{
					Initialized: true,
					Value:       4,
				},
				Children: &[]*localStateTree{
					{
						LocalStateClosure: &localStateClosure[any]{
							Initialized: true,
							Value:       V{v: 5},
						},
						Children: &[]*localStateTree{},
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
