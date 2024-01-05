package application

import (
	"github.com/BooleanCat/go-functional/iter"
	"github.com/nktknshn/go-tg-bot/tgbot/component"
	"github.com/nktknshn/go-tg-bot/tgbot/outcoming"
	"github.com/nktknshn/go-tg-bot/tgbot/render"
	"go.uber.org/zap"
)

type preRenderData[S any, C any] struct {
	NextChatState ApplicationChatState[S, C]
	RenderActions []render.RenderAction
	// ExecuteRender func(ctx context.Context, renderer ChatRenderer) ([]RenderedElement, error)
}

func (d *preRenderData[S, C]) RenderActionsKinds() []string {
	return iter.Map(iter.Lift(d.RenderActions), func(a render.RenderAction) string {
		return a.RenderActionKind()
	}).Collect()
}

type ComputeNextStateProps struct {
	Logger *zap.Logger
}

// ComputeNextState computes the output based on the state
func (app *Application[S, C]) ComputeNextState(chatState *ApplicationChatState[S, C], props ComputeNextStateProps) *preRenderData[S, C] {

	comp := app.StateToComp(chatState.AppState)

	globalContext := app.globalContext(chatState)

	createElementsResult := component.CreateElements(
		comp,
		globalContext,
		chatState.treeState,
		props.Logger,
	)

	res := outcoming.ElementsToMessagesAndHandlers(createElementsResult.Elements)

	nextChatState := ApplicationChatState[S, C]{
		ChatID:           chatState.ChatID,
		AppState:         chatState.AppState,
		renderedElements: chatState.renderedElements,
		inputHandler:     res.InputHandler,
		callbackHandler:  res.CallbackHandler,
		Renderer:         chatState.Renderer,
		treeState:        &createElementsResult.TreeState,
		lock:             chatState.lock,
	}

	return &preRenderData[S, C]{
		NextChatState: nextChatState,
		RenderActions: render.CreateRenderActions(
			chatState.renderedElements,
			res.OutcomingMessages,
			props.Logger,
		),
	}

}
