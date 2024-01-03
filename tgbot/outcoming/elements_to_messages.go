package outcoming

import (
	"fmt"

	"github.com/nktknshn/go-tg-bot/tgbot/common"
	"github.com/nktknshn/go-tg-bot/tgbot/component"
	"github.com/nktknshn/go-tg-bot/tgbot/logging"
)

type inputHandlersType []component.ElementInputHandler

type processElementsResult struct {
	OutcomingMessages []OutcomingMessage
	InputHandlers     inputHandlersType
	CallbackMap       callbackMap
	BottomButtons     []component.ElementBottomButton
	isComplete        bool

	CallbackHandler common.ChatCallbackHandler
	InputHandler    common.ChatInputHandler
}

func inputHandler(ihs inputHandlersType) common.ChatInputHandler {

	if len(ihs) == 0 {
		return nil
	}

	return func(text string) any {

		for _, h := range ihs {
			res := h.Handler(text)

			_, goNext := res.(common.ActionNext)

			if !goNext {
				return res
			}

		}
		return common.ActionNext{}
	}
}
func (per *processElementsResult) lastTextMessage() *OutcomingTextMessage {
	for i := len(per.OutcomingMessages) - 1; i >= 0; i-- {
		if per.OutcomingMessages[i].OutcomingKind() == KindOutcomingTextMessage {
			return per.OutcomingMessages[i].(*OutcomingTextMessage)
		}
	}
	return nil
}

// adds keyboard extra to the last message
func (per *processElementsResult) Complete() {

	if per.isComplete {
		logging.Logger().Error("ProcessElementsResult.Complete: already complete")
		return
	}

	lastMessage := per.lastTextMessage()

	if len(per.BottomButtons) > 0 {
		lastMessage.BottomButtons = append(lastMessage.BottomButtons, per.BottomButtons...)
	}

	if lastMessage == nil {
		logging.Logger().Error("ProcessElementsResult.Complete: no text messages")
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

func ElementsToMessagesAndHandlers(elements []component.AnyElement) *processElementsResult {
	messages := make([]OutcomingMessage, 0)
	inputHandlers := make([]component.ElementInputHandler, 0)
	bottomButtons := make([]component.ElementBottomButton, 0)
	// callbackHandlers := make(map[string]func() A)

	var lastMessage *OutcomingTextMessage

	getLastMessage := func() *OutcomingTextMessage {
		if lastMessage != nil {
			return lastMessage
		}

		for _, message := range messages {
			if message.OutcomingKind() == KindOutcomingTextMessage {
				lastMessage = message.(*OutcomingTextMessage)
			}
		}

		if lastMessage == nil {
			lastMessage = NewOutcomingTextMessage("")
			messages = append(messages, lastMessage)
		}

		return lastMessage
	}

	for _, element := range elements {
		switch el := element.(type) {

		case *component.ElementMessage:

			outcoming := NewOutcomingTextMessage(el.Text)
			lastMessage = outcoming
			messages = append(messages, outcoming)

		case *component.ElementMessagePart:
			_lastMessage := getLastMessage()

			text := element.(*component.ElementMessagePart).Text

			if _lastMessage.isComplete {
				lastMessage = NewOutcomingTextMessage(text)
				messages = append(messages, lastMessage)
			} else {
				getLastMessage().ConcatText(text)
			}

		case *component.ElementCompleteMessage:
			if lastMessage != nil {
				getLastMessage().SetComplete()
			}

		case *component.ElementButton:
			getLastMessage().AddButton(element.(*component.ElementButton))

		case *component.ElementButtonsRow:
			getLastMessage().AddButtonsRow(element.(*component.ElementButtonsRow))

		case *component.ElementInputHandler:
			inputHandlers = append(inputHandlers, *el)

		case *component.ElementBottomButton:
			bottomButtons = append(bottomButtons, *el)

		case *component.ElementUserMessage:
			messages = append(messages, &OutcomingUserMessage{*el})
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
		InputHandler:      inputHandler(inputHandlers),
	}
}
