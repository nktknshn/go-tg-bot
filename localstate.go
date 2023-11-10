package tgbot

type LocalStateClosure[S any] struct {
	value *S
}

type LocalStateTree[S any] struct {
	localState *LocalStateClosure[S]
	children   *[]LocalStateTree[any]
}

type LocalStateWithGetSet[S any] struct {
	getset     GetSetLocalStateImpl[S]
	localState *LocalStateClosure[S]
}

func NewLocalState[S any](index []int, localState *LocalStateClosure[S]) LocalStateWithGetSet[S] {

	if localState == nil {
		// local state hoder struct
		localState = &LocalStateClosure[S]{}
	}

	return LocalStateWithGetSet[S]{
		getset:     createGetSet(index, localState),
		localState: localState,
	}
}

func NewLocalStateTree[S any]() *LocalStateTree[S] {
	return &LocalStateTree[S]{
		localState: nil,
		children:   nil,
	}
}

type ActionLocalState[S any] struct {
	index []int
	f     func(S)
}

type GetSetLocalStateImpl[S any] struct {
	localState *LocalStateClosure[S]
	index      []int
}

func (g GetSetLocalStateImpl[S]) Get(initialValue S) S {

	if g.localState.value == nil {
		g.localState.value = &initialValue
	}

	return *g.localState.value
}

func (g GetSetLocalStateImpl[S]) Set(f func(S)) ActionLocalState[S] {
	return ActionLocalState[S]{
		index: g.index,
		f:     f,
	}
}

func createGetSet[S any](index []int, localState *LocalStateClosure[S]) GetSetLocalStateImpl[S] {
	return GetSetLocalStateImpl[S]{
		index:      index,
		localState: localState,
	}
}
