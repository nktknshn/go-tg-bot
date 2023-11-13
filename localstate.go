package tgbot

type LocalStateClosure[S any] struct {
	Initialized bool
	Value       S
}

// tree of the local states of components
type LocalStateTree struct {
	// local state of the current component
	LocalStateClosure *LocalStateClosure[any]
	// local states of the Children components
	// if nil then the state has to be reinitialized
	Children *[]*LocalStateTree
}

func (lst *LocalStateTree) Set(index []int, f func(any) any) {
	closure := lst.Get(index)
	closure.Value = f(closure.Value)
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
	Getset     GetSetLocalStateImpl[S]
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
	index []int
	f     func(S) S
}

type GetSetLocalStateImpl[S any] struct {
	LocalState LocalStateClosure[S]
	Index      []int
}

type GetSetStruct[S any, A any] struct {
	Get func() S
	Set func(func(S) S) A
}

func (g *GetSetLocalStateImpl[S]) Init(initialValue S) GetSetStruct[S, any] {

	if !g.LocalState.Initialized {
		globalLogger.Debug("Initializing")
		g.LocalState.Value = initialValue
		g.LocalState.Initialized = true
	}

	return GetSetStruct[S, any]{
		Get: func() S {
			return g.LocalState.Value
		},
		Set: func(f func(S) S) any {
			return ActionLocalState[S]{
				index: g.Index,
				f:     f,
			}
		},
	}
}

func NewGetSet[S any](index []int, localState *LocalStateClosure[S]) GetSetLocalStateImpl[S] {
	return GetSetLocalStateImpl[S]{
		Index:      index,
		LocalState: *localState,
	}
}
