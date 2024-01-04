package outcoming

import (
	"fmt"

	"github.com/nktknshn/go-tg-bot/tgbot/common"
	"github.com/nktknshn/go-tg-bot/tgbot/component"
	"github.com/nktknshn/go-tg-bot/tgbot/logging"
)

type outcomingResult struct {
	OutcomingMessages []OutcomingMessage
	InputHandlers     []component.ElementInputHandler
	CallbackMap       callbackMap
	BottomButtons     []component.ElementBottomButton
	isComplete        bool

	CallbackHandler common.ChatCallbackHandler
	InputHandler    common.ChatInputHandler
}

func (per *outcomingResult) lastTextMessage() *OutcomingTextMessage {
	for i := len(per.OutcomingMessages) - 1; i >= 0; i-- {
		if per.OutcomingMessages[i].OutcomingKind() == KindOutcomingTextMessage {
			return per.OutcomingMessages[i].(*OutcomingTextMessage)
		}
	}
	return nil
}

// adds keyboard extra to the last message
func (per *outcomingResult) Complete() {

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

func (per *outcomingResult) String() string {
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
