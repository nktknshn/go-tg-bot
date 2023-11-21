package tgbot

type EmptyGlobalContext struct{}

func newEmptyGlobalContext() globalContextTypedImpl[EmptyGlobalContext] {
	return newGlobalContextTyped[EmptyGlobalContext](EmptyGlobalContext{})
}

type globalContext[C any] interface {
	Get() any
}

type globalContextTypedImpl[C any] struct {
	Values C
}

func newGlobalContextTyped[C any](values C) globalContextTypedImpl[C] {
	return globalContextTypedImpl[C]{
		Values: values,
	}
}

func (c globalContextTypedImpl[C]) Get() any {
	return c.Values
}
