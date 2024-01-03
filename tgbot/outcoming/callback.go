package outcoming

import "github.com/nktknshn/go-tg-bot/tgbot/common"

type callbackMap map[string]func() *common.CallbackResult

func getCallbackHandlersMap(outcomingMessages []OutcomingMessage) callbackMap {

	callbackHandlers := make(map[string]func() *common.CallbackResult)

	for _, m := range outcomingMessages {
		switch el := m.(type) {
		case *OutcomingTextMessage:
			for _, row := range el.Buttons {
				for _, butt := range row {

					butt := butt
					callbackHandlers[butt.CallbackData()] = func() *common.CallbackResult {
						v := butt.OnClick()

						return &common.CallbackResult{
							Action:   v,
							NoAnswer: butt.NoCallback,
						}
					}
				}
			}
		}
	}

	return callbackMap(callbackHandlers)
}

func callbackMapToHandler(cbmap callbackMap) common.ChatCallbackHandler {
	return func(callbackData string) *common.CallbackResult {

		if handler, ok := cbmap[callbackData]; ok {

			return handler()
		} else {
			return nil
		}

	}
}
