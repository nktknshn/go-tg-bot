package tgbot

import (
	"fmt"

	"go.uber.org/zap"
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
	Button(string, func() A)
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

func (o *outputImpl[A]) Button(text string, handler func() A) {
	o.Result = append(o.Result, Button(text, handler, text, false))
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

		globalLogger.Info("Callback handler", zap.String("data", callbackData))

		if handler, ok := cbmap[callbackData]; ok {
			globalLogger.Info("Calling handler", zap.String("data", callbackData))

			return handler()
		} else {
			globalLogger.Error("No handler for callback", zap.String("key", callbackData))
			return nil
		}

	}
}

type ProcessElementsResult[A any] struct {
	OutcomingMessages []OutcomingMessage
	InputHandlers     []ElementInputHandler[A]
	CallbackHandler   ChatCallbackHandler[A]
	CallbackMap       map[string]func() *A
	BottomButtons     []ElementBottomButton
	isComplete        bool
}

func (per *ProcessElementsResult[A]) lastTextMessage() *OutcomingTextMessage[A] {
	for i := len(per.OutcomingMessages) - 1; i >= 0; i-- {
		if per.OutcomingMessages[i].OutcomingKind() == KindOutcomingTextMessage {
			return per.OutcomingMessages[i].(*OutcomingTextMessage[A])
		}
	}
	return nil
}

// adds keyboard extra to the last message
func (per *ProcessElementsResult[A]) Complete() {

	if per.isComplete {
		globalLogger.Error("ProcessElementsResult.Complete: already complete")
		return
	}

	lastMessage := per.lastTextMessage()

	if len(per.BottomButtons) > 0 {
		lastMessage.BottomButtons = append(lastMessage.BottomButtons, per.BottomButtons...)
	}

	if lastMessage == nil {
		globalLogger.Error("ProcessElementsResult.Complete: no text messages")
		return
	}

	per.isComplete = true
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
	bottomButtons := make([]ElementBottomButton, 0)
	// callbackHandlers := make(map[string]func() A)

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
			_lastMessage := getLastMessage()
			text := element.(*ElementMessagePart).Text

			if _lastMessage.isComplete {
				messages = append(messages, NewOutcomingTextMessage[A](text))
			} else {
				getLastMessage().ConcatText(text)

			}

		case *ElementCompleteMessage:
			getLastMessage().SetComplete()

		case *ElementButton[A]:
			getLastMessage().AddButton(element.(*ElementButton[A]))

			// TODO create callback handler

		case *ElementButtonsRow[A]:
			getLastMessage().AddButtonsRow(element.(*ElementButtonsRow[A]))

		case *ElementInputHandler[A]:
			inputHandlers = append(inputHandlers, *el)

		case *ElementBottomButton:
			bottomButtons = append(bottomButtons, *el)

		case *ElementUserMessage:
			messages = append(messages, &OutcomingUserMessage{*el})
		}
	}

	callbackMap := getCallbackHandlersMap[A](messages)
	callbackHandler := callbackMapToHandler[A](callbackMap)

	return &ProcessElementsResult[A]{
		OutcomingMessages: messages,
		InputHandlers:     inputHandlers,
		CallbackMap:       callbackMap,
		CallbackHandler:   callbackHandler,
		BottomButtons:     bottomButtons,
	}
}
