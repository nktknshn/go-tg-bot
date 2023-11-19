package tgbot

import (
	"fmt"
)

type LocalStateSetter[S any] interface {
	Set(func(S)) S
}

type LocalStateGetter[S any] interface {
	Get() S
}

// type LocalStateIniter[S any] interface {
// 	Init(S) GetSetLocalState[S, any]
// }

type GetSetLocalState[S any] interface {
	LocalStateSetter[S]
	LocalStateGetter[S]
}

type LocalStateProvider[S any] interface {
	GetSetLocalState[S]
}

type Z interface {
	O[any]
	CompZ(func(any, Z), any)
}

type O[A any] interface {
	Send(Element)
	Comp(Comp[A])
	Message(string)
	Messagef(string, ...interface{})
	MessagePart(string)
	MessagePartf(string, ...interface{})
	Button(string, func() A, ...interface{})
	ButtonsRow([]string, func(int, string) A)
	BottomButton(string)
	MessageComplete()
	InputHandler(func(string) any)
	// Dispatch(A)
	// LocalStateProvider[any, A]
}

type OO = O[any]

// type Comp[A any] func(O[A])
// type ComponentCons[T any, A any] func(props *T) Comp[A]

type Comp[A any] interface {
	Render(O[A])
}

type outputImpl[A any] struct {
	Result []Element
}

type NoCallbackStruct struct{}
type NextRowStruct struct{}

var NoCallback = NoCallbackStruct{}
var NextRow = NextRowStruct{}

func NewOutput[A any]() *outputImpl[A] {
	return &outputImpl[A]{Result: make([]Element, 0)}
}

func (o *outputImpl[A]) Message(text string) {
	o.Result = append(o.Result, Message(text))
}

func (o *outputImpl[A]) Messagef(format string, args ...interface{}) {
	o.Result = append(o.Result, Message(fmt.Sprintf(format, args...)))
}

func (o *outputImpl[A]) MessagePart(text string) {
	o.Result = append(o.Result, MessagePart(text))
}

func (o *outputImpl[A]) MessagePartf(format string, args ...interface{}) {
	o.Result = append(o.Result, MessagePart(fmt.Sprintf(format, args...)))
}

func (o *outputImpl[A]) Button(text string, handler func() A, options ...interface{}) {
	var (
		noCallback = false
		nextRow    = false
	)

	if options == nil {
		options = make([]interface{}, 0)
	}

	for _, option := range options {
		switch option.(type) {
		case NoCallbackStruct:
			noCallback = true
		case NextRowStruct:
			nextRow = true
		}
	}

	o.Result = append(o.Result, Button(text, handler, text, nextRow, noCallback))
}

func (o *outputImpl[A]) ButtonsRow(texts []string, handler func(int, string) A) {
	o.Result = append(o.Result, ButtonsRow(texts, handler))
}

func (o *outputImpl[A]) BottomButton(text string) {
	o.Result = append(o.Result, MessagePart(text))
}

func (o *outputImpl[A]) Send(element Element) {
	o.Result = append(o.Result, element)
}

func (o *outputImpl[A]) Comp(comp Comp[A]) {
	o.Result = append(o.Result, Component(comp))
}

func (o *outputImpl[A]) MessageComplete() {
	o.Result = append(o.Result, MessageComplete())
}

func (o *outputImpl[A]) InputHandler(handler func(string) any) {
	o.Result = append(o.Result, AInputHandler[A](handler))
}
