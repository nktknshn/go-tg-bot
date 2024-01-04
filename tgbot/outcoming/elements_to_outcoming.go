package outcoming

import (
	"github.com/nktknshn/go-tg-bot/tgbot/component"
)

func ElementsToMessagesAndHandlers(elements []component.AnyElement) *outcomingResult {
	messages := make([]OutcomingMessage, 0)
	inputHandlers := make([]component.ElementInputHandler, 0)
	bottomButtons := make([]component.ElementBottomButton, 0)

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

	return &outcomingResult{
		OutcomingMessages: messages,
		InputHandlers:     inputHandlers,
		CallbackMap:       callbackMap,
		CallbackHandler:   callbackHandler,
		BottomButtons:     bottomButtons,
		InputHandler:      buildInputHandler(inputHandlers),
	}
}
