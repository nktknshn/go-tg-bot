package tgbot

import (
	"context"

	"go.uber.org/zap"
)

type preRenderData[S any, C any] struct {
	InternalChatState ChatState[S, C]
	ExecuteRender     func(renderer ChatRenderer) ([]RenderedElement, error)
}

func (a *Application[S, C]) PreRender(ac *ApplicationContext[S, C]) *preRenderData[S, C] {
	logger := zap.NewNop()

	comp := ac.App.StateToComp(ac.State.AppState)

	var globalContext globalContext[C]

	if ac.App.CreateGlobalContext != nil {
		ctxValue := ac.App.CreateGlobalContext(ac.State)
		globalContext = newGlobalContextTyped[C](ctxValue)
	} else {
		globalContext = newEmptyGlobalContext()
	}

	createElementsResult := createElements(
		comp,
		globalContext,
		ac.State.treeState,
	)

	els := createElementsResult.Elements

	res := elementsToMessagesAndHandlers(els)

	logger.Debug("PreRender",
		zap.Any("OutcomingMessages", res))

	var inputHandler chatInputHandler

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
					zap.Any("res", reflectStructName(res)),
				)

				_, goNext := res.(ActionNext)

				if !goNext {
					return res
				}

			}
			return ActionNext{}
		}
	}

	nextState := ChatState[S, C]{
		ChatID:           ac.State.ChatID,
		AppState:         ac.State.AppState,
		renderedElements: ac.State.renderedElements,
		inputHandler:     inputHandler,
		callbackHandler:  res.CallbackHandler,
		Renderer:         ac.State.Renderer,
		treeState:        &createElementsResult.TreeState,
		lock:             ac.State.lock,
	}

	return &preRenderData[S, C]{
		InternalChatState: nextState,
		ExecuteRender: func(renderer ChatRenderer) ([]RenderedElement, error) {
			logger.Info("ExecuteRender")

			actions := getRenderActions(
				ac.State.renderedElements,
				res.OutcomingMessages,
			)

			for _, a := range actions {
				logger.Info("RenderActions", zap.Any("action", a))
			}

			rendered, err := executeRenderActions(
				context.Background(),
				ac.State.Renderer,
				actions,
			)

			logger.Info("Rendered", zap.Any("count", len(rendered)))

			if err != nil {
				logger.Error("Error rendering", zap.Error(err))
				return nil, err
			}

			for _, r := range rendered {
				logger.Info("Rendered", zap.Any("element", r.renderedKind()))
			}

			return rendered, nil

		},
	}

}
