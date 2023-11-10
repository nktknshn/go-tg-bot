package tgbot

import (
	"context"

	"go.uber.org/zap"
)

type PreRenderData[S any, A any] struct {
	InternalChatState InternalChatState[S, A]
	ExecuteRender     func(renderer ChatRenderer) ([]RenderedElement, error)
}

func (a *Application[S, A]) PreRender(ac *ApplicationContext[S, A]) *PreRenderData[S, A] {

	els := ComponentToElements(
		ac.App.StateToComp(ac.State.AppState),
		ac.Logger,
	)

	res := ElementsToMessagesAndHandlers[A](els)

	ac.Logger.Debug("PreRender", zap.Any("elements", res))

	var inputHandler ChatInputHandler[A]

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
		CallbackHandler:  res.CallbackHandler,
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

			for _, a := range actions {
				ac.Logger.Info("RenderActions", zap.Any("action", a))
			}

			rendered, err := ExecuteRenderActions[A](
				context.Background(),
				ac.State.Renderer,
				actions,
			)

			ac.Logger.Info("Rendered", zap.Any("count", len(rendered)))

			if err != nil {
				ac.Logger.Error("Error rendering", zap.Error(err))
				return []RenderedElement{}, err
			}

			for _, r := range rendered {
				ac.Logger.Info("Rendered", zap.Any("element", r.renderedKind()))
			}

			return rendered, nil

		},
	}

}
