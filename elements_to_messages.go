package tgbot

import (
	"fmt"

	"go.uber.org/zap"
)

type processElementsResult struct {
	OutcomingMessages []outcomingMessage
	InputHandlers     []elementInputHandler
	CallbackHandler   chatCallbackHandler
	CallbackMap       callbackMap
	BottomButtons     []elementBottomButton
	isComplete        bool
}

func (per *processElementsResult) lastTextMessage() *outcomingTextMessage {
	for i := len(per.OutcomingMessages) - 1; i >= 0; i-- {
		if per.OutcomingMessages[i].OutcomingKind() == kindOutcomingTextMessage {
			return per.OutcomingMessages[i].(*outcomingTextMessage)
		}
	}
	return nil
}

// adds keyboard extra to the last message
func (per *processElementsResult) Complete() {

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

func (per *processElementsResult) String() string {
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

func elementsToMessagesAndHandlers(elements []anyElement) *processElementsResult {
	messages := make([]outcomingMessage, 0)
	inputHandlers := make([]elementInputHandler, 0)
	bottomButtons := make([]elementBottomButton, 0)
	// callbackHandlers := make(map[string]func() A)

	var lastMessage *outcomingTextMessage

	getLastMessage := func() *outcomingTextMessage {
		if lastMessage != nil {
			return lastMessage
		}

		for _, message := range messages {
			if message.OutcomingKind() == kindOutcomingTextMessage {
				lastMessage = message.(*outcomingTextMessage)
			}
		}

		if lastMessage == nil {
			lastMessage = newOutcomingTextMessage("")
			messages = append(messages, lastMessage)
		}

		return lastMessage
	}

	for _, element := range elements {
		switch el := element.(type) {

		case *elementMessage:
			globalLogger.Debug("ElementMessage", zap.Any("el", el))

			outcoming := newOutcomingTextMessage(el.Text)
			lastMessage = outcoming
			messages = append(messages, outcoming)

		case *elementMessagePart:
			_lastMessage := getLastMessage()

			text := element.(*elementMessagePart).Text

			if _lastMessage.isComplete {
				lastMessage = newOutcomingTextMessage(text)
				messages = append(messages, lastMessage)
			} else {
				getLastMessage().ConcatText(text)
			}

		case *elementCompleteMessage:
			if lastMessage != nil {
				getLastMessage().SetComplete()
			}

		case *elementButton:
			getLastMessage().AddButton(element.(*elementButton))

		case *elementButtonsRow:
			getLastMessage().AddButtonsRow(element.(*elementButtonsRow))

		case *elementInputHandler:
			inputHandlers = append(inputHandlers, *el)

		case *elementBottomButton:
			bottomButtons = append(bottomButtons, *el)

		case *elementUserMessage:
			messages = append(messages, &outcomingUserMessage{*el})
		}
	}

	callbackMap := getCallbackHandlersMap(messages)
	callbackHandler := callbackMapToHandler(callbackMap)

	return &processElementsResult{
		OutcomingMessages: messages,
		InputHandlers:     inputHandlers,
		CallbackMap:       callbackMap,
		CallbackHandler:   callbackHandler,
		BottomButtons:     bottomButtons,
	}
}
