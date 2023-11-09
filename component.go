package tgbot

import (
	"fmt"

	"go.uber.org/zap"
)

type O[A any] interface {
	Send(Element)
	Comp(Comp[A])
	Message(string)
	Messagef(string, ...interface{})
	MessagePart(string)
	Button(string, func() A)
	ButtonsRow([]string, func(int, string) A)
	BottomButton(string)
	MessageComplete()
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
	o.result = append(o.result, Button(text, handler, text, false))
}

func (o *outputImpl[A]) ButtonsRow(texts []string, handler func(int, string) A) {
	o.result = append(o.result, ButtonsRow(texts, handler))
}

func (o *outputImpl[A]) BottomButton(text string) {
	o.result = append(o.result, MessagePart(text))
}

func (o *outputImpl[A]) Send(element Element) {
	o.result = append(o.result, element)
}

func (o *outputImpl[A]) Comp(comp Comp[A]) {
	o.result = append(o.result, Component(comp))
}

func (o *outputImpl[A]) MessageComplete() {
	o.result = append(o.result, MessageComplete())
}

func (o *outputImpl[A]) InputHandler(handler func(string) A) {
	o.result = append(o.result, AInputHandler(handler))
}

func ComponentToElements[A any](comp Comp[A]) []Element {
	o := &outputImpl[A]{result: make([]Element, 0)}
	comp(o)
	return o.result
}

func getCallbackHandlersMap[A any](outcomingMessages []OutcomingMessage) map[string]func() *A {

	callbackHandlers := make(map[string]func() *A)

	for _, m := range outcomingMessages {
		switch el := m.(type) {
		case *OutcomingTextMessage[A]:
			for _, row := range el.Buttons {
				for _, butt := range row {

					butt := butt
					callbackHandlers[butt.CallbackData()] = func() *A {
						v := butt.OnClick()
						return &v
					}
				}
			}
		}
	}

	return callbackHandlers
}

func callbackMapToHandler[A any](cbmap map[string]func() *A) ChatCallbackHandler[A] {
	return func(callbackData string) *A {

		logger.Info("Callback handler", zap.String("data", callbackData))

		if handler, ok := cbmap[callbackData]; ok {
			logger.Info("Calling handler", zap.String("data", callbackData))

			return handler()
		} else {
			logger.Error("No handler for callback", zap.String("key", callbackData))
			return nil
		}

	}
}

type ProcessElementsResult[A any] struct {
	OutcomingMessages []OutcomingMessage
	InputHandlers     []ElementInputHandler[A]
	CallbackHandler   ChatCallbackHandler[A]
	CallbackMap       map[string]func() *A
}

func (per *ProcessElementsResult[A]) String() string {
	messagesStr := ""

	for _, m := range per.OutcomingMessages {
		messagesStr += fmt.Sprintf("%v,", m.OutcomingKind())
	}

	cbhsStr := ""

	for k := range per.CallbackMap {
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

	callbackMap := getCallbackHandlersMap[A](messages)
	callbackHandler := callbackMapToHandler[A](callbackMap)

	return &ProcessElementsResult[A]{
		OutcomingMessages: messages,
		InputHandlers:     inputHandlers,
		CallbackMap:       callbackMap,
		CallbackHandler:   callbackHandler,
	}
}
