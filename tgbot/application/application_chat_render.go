package application

import (
	"context"

	"github.com/nktknshn/go-tg-bot/tgbot/render"
	"go.uber.org/zap"
)

// Computes the output based on the state and renders it to the user
func DefaultRenderFunc[S any, C any](ctx context.Context, ac *ApplicationChat[S, C]) error {

	logger := ac.Loggers.Render()

	logger.Debug("RenderFunc. Computing next state")

	res := ac.App.ComputeNextState(ac.State, ComputeNextStateProps{
		Logger: ac.Loggers.Component(),
	})

	logger.Debug("RenderFunc computed next state",
		zap.Any("RenderActions", res.RenderActionsKinds()),
	)

	rendered, err := render.ExecuteRenderActions(
		ctx,
		ac.State.Renderer,
		res.RenderActions,
		render.ExecuteRenderActionsProps{Logger: ac.Loggers.Render()})

	if err != nil {
		logger.Error("Error in RenderFunc", zap.Error(err))
		return err
	}

	ac.SetChatState(&res.NextChatState)
	ac.State.SetRenderedElements(rendered)

	return nil
}
