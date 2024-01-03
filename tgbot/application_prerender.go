package tgbot

import (
	"github.com/BooleanCat/go-functional/iter"
	"go.uber.org/zap"
)

type preRenderData[S any, C any] struct {
	NextChatState ChatState[S, C]
	RenderActions []renderActionType
	// ExecuteRender func(ctx context.Context, renderer ChatRenderer) ([]RenderedElement, error)
}

func (d *preRenderData[S, C]) RenderActionsKinds() []string {
	return iter.Map(iter.Lift(d.RenderActions), func(a renderActionType) string {
		return a.RenderActionKind()
	}).Collect()
}

// ComputeNextState computes the output based on the state
func (app *Application[S, C]) ComputeNextState(chatState *ChatState[S, C], logger *zap.Logger) *preRenderData[S, C] {

	comp := app.StateToComp(chatState.AppState)

	globalContext := app.globalContext(chatState)

	createElementsResult := createElements(
		comp,
		globalContext,
		chatState.treeState,
		logger,
	)

	res := elementsToMessagesAndHandlers(createElementsResult.Elements)

	nextChatState := ChatState[S, C]{
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
		RenderActions: getRenderActions(
			chatState.renderedElements,
			res.OutcomingMessages,
			logger,
		),
	}

}
