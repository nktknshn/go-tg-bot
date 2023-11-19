package tgbot

import (
	"fmt"

	"go.uber.org/zap"
)

type ProcessElementsResult[A any] struct {
	OutcomingMessages []OutcomingMessage
	InputHandlers     []ElementInputHandler[A]
	CallbackHandler   ChatCallbackHandler[A]
	CallbackMap       callbackMap[A]
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
		messagesStr += fmt.Sprintf("%v,", m)
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
			globalLogger.Debug("ElementMessage", zap.Any("el", el))

			outcoming := NewOutcomingTextMessage[A](el.Text)
			lastMessage = outcoming
			messages = append(messages, outcoming)

		case *ElementMessagePart:
			_lastMessage := getLastMessage()

			text := element.(*ElementMessagePart).Text

			if _lastMessage.isComplete {
				lastMessage = NewOutcomingTextMessage[A](text)
				messages = append(messages, lastMessage)
			} else {
				getLastMessage().ConcatText(text)
			}

		case *ElementCompleteMessage:
			if lastMessage != nil {
				getLastMessage().SetComplete()
			}

		case *ElementButton[A]:
			getLastMessage().AddButton(element.(*ElementButton[A]))

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
