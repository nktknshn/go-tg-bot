package component

import (
	"fmt"
)

type localStateClosure[S any] struct {
	Initialized bool
	Value       S
}

func (lsc *localStateClosure[S]) String() string {
	return fmt.Sprintf("Initialized: %v, Value: %v", lsc.Initialized, lsc.Value)
}

// tree of the local states of components
type localStateTree struct {
	CompId string
	// local state of the current component
	LocalStateClosure *localStateClosure[any]
	// local states of the Children components
	// if nil then the state has to be reinitialized
	Children *[]*localStateTree
}

func (lst *localStateTree) String() string {
	result := ""

	result += fmt.Sprintf("CompId: %v, ", lst.CompId)
	result += fmt.Sprintf("LocalStateClosure: {%v}, ", lst.LocalStateClosure)

	if lst.Children == nil {
		return result + "Children: nil"
	}

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

func (lst *localStateTree) Set(index []int, f func(any) any) {
	closure := lst.Get(index)
	closure.Value = f(closure.Value)
	closure.Initialized = true
}

func (lst *localStateTree) Get(index []int) *localStateClosure[any] {
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

type localStateWithGetSet[S any] struct {
	Getset     CompState[S]
	LocalState *localStateClosure[S]
}

// creates an empty closure for local state and get and set functions
func newLocalState[S any](index []int, localState *localStateClosure[S]) localStateWithGetSet[S] {

	if localState == nil {
		// local state hoder struct
		localState = &localStateClosure[S]{}
	}

	return localStateWithGetSet[S]{
		Getset:     newGetSet(index, localState),
		LocalState: localState,
	}
}

func newLocalStateTree() *localStateTree {
	return &localStateTree{
		LocalStateClosure: nil,
		Children:          nil,
	}
}

type ActionLocalState[S any] struct {
	Index []int
	F     func(S) S
}

type CompState[S any] struct {
	LocalStateClosure *localStateClosure[S]
	Index             []int
}

type getSetStruct[S any, A any] struct {
	Get func() S
	Set func(func(S) S) A
}

func (g CompState[S]) Init(initialValue S) getSetStruct[S, any] {

	if !g.LocalStateClosure.Initialized {
		// globalLogger.Debug("Initializing",
		// 	zap.Any("index", g.Index),
		// 	zap.Any("initialValue", initialValue),
		// )

		g.LocalStateClosure.Value = initialValue
		g.LocalStateClosure.Initialized = true
	}

	return getSetStruct[S, any]{
		Get: func() S {
			return g.LocalStateClosure.Value
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

func newGetSet[S any](index []int, localState *localStateClosure[S]) CompState[S] {
	return CompState[S]{
		Index:             index,
		LocalStateClosure: localState,
	}
}
