package tgbot

import (
	"fmt"
)

type O[A any] interface {
	Send(Element)
	Comp(Comp[A])
	Message(string)
	Messagef(string, ...interface{})
	MessagePart(string)
	Button(string, func() A)
	EndMessage()
	InputHandler(func(string) A)
}

type Comp[A any] func(O[A])
type ComponentCons[T any, A any] func(props *T) Comp[A]

type outputImpl[A any] struct {
	result []Element
}

func (o *outputImpl[A]) Message(text string) {
	o.result = append(o.result, Message(text))
}

func (o *outputImpl[A]) Messagef(format string, args ...interface{}) {
	o.result = append(o.result, Message(fmt.Sprintf(format, args...)))
}

func (o *outputImpl[A]) MessagePart(text string) {
	o.result = append(o.result, MessagePart(text))
}

func (o *outputImpl[A]) Button(text string, handler func() A) {
	o.result = append(o.result, Button(text, handler))
}

func (o *outputImpl[A]) Send(element Element) {
	o.result = append(o.result, element)
}

func (o *outputImpl[A]) Comp(comp Comp[A]) {
	o.result = append(o.result, Component(comp))
}

func (o *outputImpl[A]) EndMessage() {
	o.result = append(o.result, EndMessage())
}

func (o *outputImpl[A]) InputHandler(handler func(string) A) {
	o.result = append(o.result, AInputHandler(handler))
}

func ComponentToElements2[A any](comp Comp[A]) []Element {

	o := &outputImpl[A]{result: make([]Element, 0)}

	comp(o)

	return o.result
}

func ComponentToElements[T any, A any](compFunc ComponentCons[T, A], props *T) []Element {

	c := compFunc(props)

	o := &outputImpl[A]{result: make([]Element, 0)}

	c(o)

	return o.result
}

type ProcessElementsResult[A any] struct {
	OutcomingMessages []OutcomingMessage
	InputHandlers     []ElementInputHandler[A]
	CallbackHandlers  map[string]ChatCallbackHandler[A]
}

func (per *ProcessElementsResult[A]) String() string {
	messagesStr := ""

	for _, m := range per.OutcomingMessages {
		messagesStr += fmt.Sprintf("%v,", m.OutcomingKind())
	}

	cbhsStr := ""

	for k := range per.CallbackHandlers {
		cbhsStr += fmt.Sprintf("%v,", k)
	}

	return fmt.Sprintf(
		"OutcomingMessages: %v, InputHandlers: %v, CallbackHandlers: %v",
		messagesStr, len(per.InputHandlers), cbhsStr,
	)
}

func ElementsToMessagesAndHandlers[A any](elements []Element) *ProcessElementsResult[A] {
	messages := make([]OutcomingMessage, 0)
	inputHandlers := make([]ElementInputHandler[A], 0)
	callbackHandlers := make(map[string]func() A)

	var lastMessage *OutcomingTextMessage[A]

	getLastMessage := func() *OutcomingTextMessage[A] {
		if lastMessage != nil {
			return lastMessage
		}

		for _, message := range messages {
			if message.OutcomingKind() == KindOutcomingTextMessage {
				lastMessage = message.(*OutcomingTextMessage[A])
			}
		}

		if lastMessage == nil {
			lastMessage = NewOutcomingTextMessage[A]("")
			messages = append(messages, lastMessage)
		}

		return lastMessage
	}

	for _, element := range elements {
		switch el := element.(type) {

		case *ElementMessage:
			outcoming := NewOutcomingTextMessage[A](el.Text)
			messages = append(messages, outcoming)

		case *ElementMessagePart:
			getLastMessage().concatText(element.(*ElementMessagePart).Text)

		case *ElementButton[A]:
			getLastMessage().AddButton(element.(*ElementButton[A]))

			act := el.Action
			if act == "" {
				act = el.Text
			}

			callbackHandlers[act] = el.OnClick

		case *ElementInputHandler[A]:
			inputHandlers = append(inputHandlers, *el)
		}
	}

	return &ProcessElementsResult[A]{
		OutcomingMessages: messages,
		InputHandlers:     inputHandlers,
	}
}
