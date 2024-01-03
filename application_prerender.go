package tgbot

import (
	"context"

	"go.uber.org/zap"
)

type preRenderData[S any, C any] struct {
	InternalChatState ChatState[S, C]
	ExecuteRender     func(renderer ChatRenderer) ([]RenderedElement, error)
}

func (app *Application[S, C]) PreRender(chatState *ChatState[S, C]) *preRenderData[S, C] {
	logger := zap.NewNop()

	comp := app.StateToComp(chatState.AppState)

	globalContext := app.globalContext(chatState)

	createElementsResult := createElements(
		comp,
		globalContext,
		chatState.treeState,
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

	nextChatState := ChatState[S, C]{
		ChatID:           chatState.ChatID,
		AppState:         chatState.AppState,
		renderedElements: chatState.renderedElements,
		inputHandler:     inputHandler,
		callbackHandler:  res.CallbackHandler,
		Renderer:         chatState.Renderer,
		treeState:        &createElementsResult.TreeState,
		lock:             chatState.lock,
	}

	return &preRenderData[S, C]{
		InternalChatState: nextChatState,
		ExecuteRender: func(renderer ChatRenderer) ([]RenderedElement, error) {
			logger.Info("ExecuteRender")

			actions := getRenderActions(
				chatState.renderedElements,
				res.OutcomingMessages,
			)

			for _, a := range actions {
				logger.Info("RenderActions", zap.Any("action", a))
			}

			rendered, err := executeRenderActions(
				context.Background(),
				chatState.Renderer,
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
