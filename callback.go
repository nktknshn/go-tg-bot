package tgbot

type callbackMap map[string]func() *callbackResult

func getCallbackHandlersMap(outcomingMessages []outcomingMessage) callbackMap {

	callbackHandlers := make(map[string]func() *callbackResult)

	for _, m := range outcomingMessages {
		switch el := m.(type) {
		case *outcomingTextMessage:
			for _, row := range el.Buttons {
				for _, butt := range row {

					butt := butt
					callbackHandlers[butt.CallbackData()] = func() *callbackResult {
						v := butt.OnClick()

						return &callbackResult{
							action:   v,
							noAnswer: butt.NoCallback,
						}
					}
				}
			}
		}
	}

	return callbackMap(callbackHandlers)
}

func callbackMapToHandler(cbmap callbackMap) chatCallbackHandler {
	return func(callbackData string) *callbackResult {

		if handler, ok := cbmap[callbackData]; ok {

			return handler()
		} else {
			return nil
		}

	}
}
