package tgbot

import (
	"fmt"

	"go.uber.org/zap"
)

type LocalStateClosure[S any] struct {
	Initialized bool
	Value       S
}

func (lsc *LocalStateClosure[S]) String() string {
	return fmt.Sprintf("Initialized: %v, Value: %v", lsc.Initialized, lsc.Value)
}

// tree of the local states of components
type LocalStateTree struct {
	CompId string
	// local state of the current component
	LocalStateClosure *LocalStateClosure[any]
	// local states of the Children components
	// if nil then the state has to be reinitialized
	Children *[]*LocalStateTree
}

func (lst *LocalStateTree) String() string {
	result := ""

	result += fmt.Sprintf("CompId: %v, ", lst.CompId)
	result += fmt.Sprintf("LocalStateClosure: {%v}, ", lst.LocalStateClosure)

	childrenStr := ""

	for _, c := range *lst.Children {
		if c == nil {
			childrenStr += "elem,"
			continue
		}
		childrenStr += fmt.Sprintf("{%v},", c)
	}

	result += fmt.Sprintf("Children: [%v]", childrenStr)

	return result
}

func (lst *LocalStateTree) Set(index []int, f func(any) any) {
	closure := lst.Get(index)
	closure.Value = f(closure.Value)
	closure.Initialized = true
}

func (lst *LocalStateTree) Get(index []int) *LocalStateClosure[any] {
	if len(index) == 0 {
		return lst.LocalStateClosure
	}

	if lst.Children == nil {
		panic("LocalStateTree.Get: Children is nil")
	}

	if len(*lst.Children) == 0 {
		panic("LocalStateTree.Get: Children is empty")
	}

	idx := index[0]

	if idx > len(*lst.Children)-1 {
		panic("LocalStateTree.Get: index out of range")
	}

	cs := (*lst.Children)[idx].Get(index[1:])

	return cs
}

type LocalStateWithGetSet[S any] struct {
	Getset     State[S]
	LocalState *LocalStateClosure[S]
}

// creates an empty closure for local state and get and set functions
func NewLocalState[S any](index []int, localState *LocalStateClosure[S]) LocalStateWithGetSet[S] {

	if localState == nil {
		// local state hoder struct
		localState = &LocalStateClosure[S]{}
	}

	return LocalStateWithGetSet[S]{
		Getset:     NewGetSet(index, localState),
		LocalState: localState,
	}
}

func NewLocalStateTree() *LocalStateTree {
	return &LocalStateTree{
		LocalStateClosure: nil,
		Children:          nil,
	}
}

type ActionLocalState[S any] struct {
	Index []int
	F     func(S) S
}

type State[S any] struct {
	LocalState LocalStateClosure[S]
	Index      []int
}

type GetSetStruct[S any, A any] struct {
	Get func() S
	Set func(func(S) S) A
}

func (g State[S]) Init(initialValue S) GetSetStruct[S, any] {

	if !g.LocalState.Initialized {
		globalLogger.Debug("Initializing", zap.Any("index", g.Index))

		g.LocalState.Value = initialValue
		g.LocalState.Initialized = true
	}

	return GetSetStruct[S, any]{
		Get: func() S {
			return g.LocalState.Value
		},
		Set: func(f func(S) S) any {
			return ActionLocalState[any]{
				Index: g.Index,
				F: func(a any) any {
					return f(a.(S))
				},
			}
		},
	}
}

func NewGetSet[S any](index []int, localState *LocalStateClosure[S]) State[S] {
	return State[S]{
		Index:      index,
		LocalState: *localState,
	}
}
