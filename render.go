package tgbot

import (
	"context"

	"go.uber.org/zap"
)

type PreRenderData[S any, A any, C any] struct {
	InternalChatState InternalChatState[S, A, C]
	ExecuteRender     func(renderer ChatRenderer) ([]RenderedElement, error)
}

func (a *Application[S, A, C]) PreRender(ac *ApplicationContext[S, A, C]) *PreRenderData[S, A, C] {
	// /	logger := ac.Logger
	logger := zap.NewNop()

	comp := ac.App.StateToComp(ac.State.AppState)

	globalContext := ac.App.CreateGlobalContext(ac.State)

	createElementsResult := CreateElements[A](
		comp,
		globalContext.(GlobalContextTyped[any]),
		ac.State.TreeState,
	)

	els := createElementsResult.Elements

	res := ElementsToMessagesAndHandlers[A](els)

	logger.Debug("PreRender",
		zap.Any("OutcomingMessages", res))

	var inputHandler ChatInputHandler[any]

	if len(res.InputHandlers) > 0 {
		handlers := res.InputHandlers

		inputHandler = func(text string) any {

			logger.Debug("InputHandler",
				zap.Any("text", text),
				zap.Any("handlers_count", len(handlers)))

			for idx, h := range handlers {
				res := h.Handler(text)

				logger.Debug("InputHandler",
					zap.Any("idx", idx),
					zap.Any("res", ReflectStructName(res)),
				)

				_, goNext := res.(Next)

				if !goNext {
					return res
				}

			}
			return Next{}
		}
	}

	nextState := InternalChatState[S, A, C]{
		ChatID:           ac.State.ChatID,
		AppState:         ac.State.AppState,
		RenderedElements: ac.State.RenderedElements,
		InputHandler:     inputHandler,
		CallbackHandler:  res.CallbackHandler,
		Renderer:         ac.State.Renderer,
		TreeState:        &createElementsResult.TreeState,
		Lock:             ac.State.Lock,
	}

	return &PreRenderData[S, A, C]{
		InternalChatState: nextState,
		ExecuteRender: func(renderer ChatRenderer) ([]RenderedElement, error) {
			logger.Info("ExecuteRender")

			actions := GetRenderActions[A](
				ac.State.RenderedElements,
				res.OutcomingMessages,
			)

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
