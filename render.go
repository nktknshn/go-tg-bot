package tgbot

import (
	"context"

	"go.uber.org/zap"
)

type PreRenderData[S any, A any] struct {
	InternalChatState InternalChatState[S, A]
	ExecuteRender     func(renderer ChatRenderer) ([]RenderedElement, error)
}

func getCallbackHandlersMap[A any](outcomingMessages []OutcomingMessage) map[string]func() *A {

	callbackHandlers := make(map[string]func() *A)

	for _, m := range outcomingMessages {
		switch el := m.(type) {
		case *OutcomingTextMessage[A]:
			for _, row := range el.Buttons {
				for _, butt := range row {

					logger.Info("Setting callback handler", zap.String("key", butt.CallbackData()))
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
			// ac.App.HandleAction(ac, tc, )
			logger.Info("Calling handler", zap.String("data", callbackData))

			return handler()
		} else {
			// err := fmt.Errorf("no handler for callback %v", key)
			logger.Error("No handler for callback", zap.String("key", callbackData))
			return nil
		}

	}
}

func (a *Application[S, A]) PreRender(ac *ApplicationContext[S, A]) *PreRenderData[S, A] {

	els := ComponentToElements2(
		ac.App.StateToComp(ac.State.AppState),
	)

	res := ElementsToMessagesAndHandlers[A](els)

	ac.Logger.Debug("PreRender", zap.Any("res", res))

	var inputHandler ChatInputHandler[A]
	callbackMap := getCallbackHandlersMap[A](res.OutcomingMessages)
	callbackHandler := callbackMapToHandler[A](callbackMap)

	if len(res.InputHandlers) > 0 {
		inputHandler = func(text string) A {
			return res.InputHandlers[0].Handler(text)
		}
	}

	nextState := InternalChatState[S, A]{
		ChatID:           ac.State.ChatID,
		AppState:         ac.State.AppState,
		RenderedElements: ac.State.RenderedElements,
		InputHandler:     inputHandler,
		CallbackHandler:  callbackHandler,
		Renderer:         ac.State.Renderer,
	}

	return &PreRenderData[S, A]{
		InternalChatState: nextState,
		ExecuteRender: func(renderer ChatRenderer) ([]RenderedElement, error) {
			ac.Logger.Info("ExecuteRender")

			actions := GetRenderActions[A](
				ac.State.RenderedElements,
				res.OutcomingMessages,
			)

			// logger.Info("RenderActions", zap.Any("actions", actions))

			for _, a := range actions {
				logger.Info("RenderActions", zap.Any("action", a))
			}

			rendered, err := ExecuteRenderActions[A](
				context.Background(),
				ac.State.Renderer,
				actions,
			)

			logger.Info("Rendered", zap.Any("count", len(rendered)))

			if err != nil {
				logger.Error("Error rendering", zap.Error(err))
				return []RenderedElement{}, err
			}

			for _, r := range rendered {
				logger.Info("Rendered", zap.Any("element", r.renderedKind()))
			}

			return rendered, nil

		},
	}

}
