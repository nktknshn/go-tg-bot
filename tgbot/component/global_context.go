package component

type EmptyGlobalContext struct{}

func NewEmptyGlobalContext() GlobalContextTypedImpl[EmptyGlobalContext] {
	return NewGlobalContextTyped[EmptyGlobalContext](EmptyGlobalContext{})
}

type GlobalContext[C any] interface {
	Get() any
}

type GlobalContextTypedImpl[C any] struct {
	Values C
}

func NewGlobalContextTyped[C any](values C) GlobalContextTypedImpl[C] {
	return GlobalContextTypedImpl[C]{
		Values: values,
	}
}

func (c GlobalContextTypedImpl[C]) Get() any {
	return c.Values
}
