package component

import (
	"fmt"
)

type O interface {
	Send(AnyElement)
	Comp(Comp)
	Message(string)
	Messagef(string, ...interface{})
	MessagePart(string)
	MessagePartf(string, ...interface{})
	Button(string, func() any, ...interface{})
	ButtonsRow([]string, func(int, string) any)
	BottomButton(string)
	MessageComplete()
	InputHandler(func(string) any)
}

type Comp interface {
	Render(O)
}

type outputImpl struct {
	Result []AnyElement
}

type noCallbackStruct struct{}
type nextRowStruct struct{}

var BtnNoCallback = noCallbackStruct{}
var BtnNextRow = nextRowStruct{}

func newOutput() *outputImpl {
	return &outputImpl{Result: make([]AnyElement, 0)}
}

func (o *outputImpl) Message(text string) {
	o.Result = append(o.Result, newMessage(text))
}

func (o *outputImpl) Messagef(format string, args ...interface{}) {
	o.Result = append(o.Result, newMessage(fmt.Sprintf(format, args...)))
}

func (o *outputImpl) MessagePart(text string) {
	o.Result = append(o.Result, newMessagePart(text))
}

func (o *outputImpl) MessagePartf(format string, args ...interface{}) {
	o.Result = append(o.Result, newMessagePart(fmt.Sprintf(format, args...)))
}

func (o *outputImpl) Button(text string, handler func() any, options ...interface{}) {
	var (
		noCallback = false
		nextRow    = false
	)

	if options == nil {
		options = make([]interface{}, 0)
	}

	for _, option := range options {
		switch option.(type) {
		case noCallbackStruct:
			noCallback = true
		case nextRowStruct:
			nextRow = true
		}
	}

	o.Result = append(o.Result, NewButton(text, handler, text, nextRow, noCallback))
}

func (o *outputImpl) ButtonsRow(texts []string, handler func(int, string) any) {
	o.Result = append(o.Result, newButtonsRow(texts, handler))
}

func (o *outputImpl) BottomButton(text string) {
	o.Result = append(o.Result, newMessagePart(text))
}

func (o *outputImpl) Send(element AnyElement) {
	o.Result = append(o.Result, element)
}

func (o *outputImpl) Comp(comp Comp) {
	o.Result = append(o.Result, newComponent(comp))
}

func (o *outputImpl) MessageComplete() {
	o.Result = append(o.Result, newMessageComplete())
}

func (o *outputImpl) InputHandler(handler func(string) any) {
	o.Result = append(o.Result, newInputHandler(handler))
}
