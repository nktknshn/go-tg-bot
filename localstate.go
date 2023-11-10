package tgbot

type ActionLocalState[S any] struct {
	index []int
	value S
}

type createLocalStateResult[S any] struct {
	// getset  GetSetLocalState[S]
	closure LocalState[S]
}

type GetSetLocalStateImpl[S any] struct {
	localState LocalState[S]
}

func (g GetSetLocalStateImpl[S]) LocalState(initialValue S) S {

	if g.localState.value == nil {
		g.localState.value = &initialValue
	}

	return *g.localState.value
}

func (g GetSetLocalStateImpl[S]) SetLocalState(value S) {
	// g.localState.value = &value
	// g.localState.updated = true
}

func createGetSet[S any](index []int, localState LocalState[S]) GetSetLocalStateImpl[S] {
	return GetSetLocalStateImpl[S]{}
}

func createLocalState[S any](index []int, localState *LocalState[S]) createLocalStateResult[S] {

	if localState == nil {
		localState = &LocalState[S]{}

		return createLocalStateResult[S]{
			closure: *localState,
		}

	}

	return createLocalStateResult[S]{
		closure: *localState,
	}
}
