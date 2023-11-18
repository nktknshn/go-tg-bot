package tgbot

import "go.uber.org/zap"

type CallbackResult[A any] struct {
	action     A
	noCallback bool
}

type callbackMap[A any] map[string]func() *CallbackResult[A]

func getCallbackHandlersMap[A any](outcomingMessages []OutcomingMessage) callbackMap[A] {

	callbackHandlers := make(map[string]func() *CallbackResult[A])

	for _, m := range outcomingMessages {
		switch el := m.(type) {
		case *OutcomingTextMessage[A]:
			for _, row := range el.Buttons {
				for _, butt := range row {

					butt := butt
					callbackHandlers[butt.CallbackData()] = func() *CallbackResult[A] {
						v := butt.OnClick()

						return &CallbackResult[A]{
							action:     v,
							noCallback: butt.NoCallback,
						}
					}
				}
			}
		}
	}

	return callbackMap[A](callbackHandlers)
}

func callbackMapToHandler[A any](cbmap callbackMap[A]) ChatCallbackHandler[A] {
	return func(callbackData string) *CallbackResult[A] {

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
