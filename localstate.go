package tgbot

type LocalStateClosure[S any] struct {
	Initialized bool
	Value       S
}

// tree of the local states of components
type LocalStateTree[S any] struct {
	// local state of the current component
	localStateClosure *LocalStateClosure[S]
	// local states of the children components
	// if nil then the state has to be reinitialized
	children *[]*LocalStateTree[any]
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

func NewLocalStateTree[S any]() *LocalStateTree[S] {
	return &LocalStateTree[S]{
		localStateClosure: nil,
		children:          nil,
	}
}

type ActionLocalState[S any] struct {
	index []int
	f     func(S)
}

type GetSetLocalStateImpl[S any] struct {
	LocalState LocalStateClosure[S]
	Index      []int
}

func (g GetSetLocalStateImpl[S]) Get(initialValue S) S {

	if !g.LocalState.Initialized {
		g.LocalState.Value = initialValue
		g.LocalState.Initialized = true
	}

	return g.LocalState.Value
}

func (g GetSetLocalStateImpl[S]) Set(f func(S)) ActionLocalState[S] {
	return ActionLocalState[S]{
		index: g.Index,
		f:     f,
	}
}

func NewGetSet[S any](index []int, localState *LocalStateClosure[S]) GetSetLocalStateImpl[S] {
	return GetSetLocalStateImpl[S]{
		Index:      index,
		LocalState: *localState,
	}
}
