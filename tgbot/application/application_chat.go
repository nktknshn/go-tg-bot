package application

import (
	"context"
	"sync"

	"github.com/nktknshn/go-tg-bot/tgbot/reflection"
	"github.com/nktknshn/go-tg-bot/tgbot/render"
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
	"go.uber.org/zap"
)

// tie together an application methods and a chat state
type ApplicationChat[S any, C any] struct {
	App   *Application[S, C]
	State *ChatState[S, C]

	// loggers for different parts of the app
	Loggers *ApplicationChatLoggers
}

func (ac *ApplicationChat[S, C]) SetChatState(chatState *ChatState[S, C]) {
	ac.State = chatState
}

func NewApplicationChat[S any, C any](app *Application[S, C], tc *telegram.TelegramUpdateContext) *ApplicationChat[S, C] {
	appState := app.CreateAppState(tc)

	chatState := ChatState[S, C]{
		ChatID:           tc.ChatID,
		AppState:         appState,
		renderedElements: []render.RenderedElement{},
		inputHandler:     nil,
		callbackHandler:  nil,
		treeState:        nil,
		Renderer:         app.CreateChatRenderer(tc),
		lock:             &sync.Mutex{},
	}

	rootLogger := app.Loggers.ApplicationChat(
		app.Loggers.Base,
	).With(zap.Int64("ChatID", tc.ChatID))

	loggers := &ApplicationChatLoggers{
		Root: rootLogger,
		// logger for components rendering
		// Component: app.Loggers.Component(rootLogger),
		Component: app.Loggers.Component(rootLogger),
		// logger for updates
		Handle: rootLogger.Named("Handle"),
		Action: rootLogger.Named("Action"),
		Render: rootLogger.Named("Render"),
	}

	res := app.ComputeNextState(&chatState, loggers.Component)

	return &ApplicationChat[S, C]{
		App:     app,
		State:   &res.NextChatState,
		Loggers: loggers,
	}
}

// Computes the output based on the state and renders it to the user
func DefaultRenderFunc[S any, C any](ctx context.Context, ac *ApplicationChat[S, C]) error {

	logger := ac.Loggers.Render

	logger.Debug("RenderFunc called")

	res := ac.App.ComputeNextState(ac.State, ac.App.Loggers.Component(logger))

	logger.Debug("RenderFunc computed next state",
		zap.Any("RenderActions", res.RenderActionsKinds()),
	)

	rendered, err := render.ExecuteRenderActions(
		ctx,
		ac.State.Renderer,
		res.RenderActions,
		render.ExecuteRenderActionsProps{Logger: logger})

	if err != nil {
		logger.Error("Error in RenderFunc", zap.Error(err))
		return err
	}

	ac.SetChatState(&res.NextChatState)
	ac.State.SetRenderedElements(rendered)

	return nil
}

func DefaultHandlerCallback[S any, C any](ac *ApplicationChat[S, C], tc *telegram.TelegramContextCallback) {

	logger := ac.Loggers.Handle.With(zap.Int64("UpdateID", tc.UpdateID))

	logger.Info("HandleCallback", zap.Any("data", tc.UpdateBotCallbackQuery.QueryID))

	// ac.Loggers.Root.Debug("LocalStateTree", zap.String("tree", ac.State.treeState.LocalStateTree.String()))

	ac.State.LockState(logger.Named("LockState"))
	defer ac.State.UnlockState(logger.Named("LockState"))

	if ac.State.callbackHandler != nil {
		result := ac.State.callbackHandler(string(tc.UpdateBotCallbackQuery.Data))

		// logger.Debug("HandleCallback", zap.Any("action", result))

		if result == nil {
			logger.Warn("CallbackHandler returned nil")
			return
		}

		internalActionHandle(ac, &tc.TelegramUpdateContext, result.Action)

		if !result.NoAnswer {
			tc.AnswerCallbackQuery()
		}

	} else {
		logger.Warn("Missing CallbackHandler")
	}

	err := ac.App.RenderFunc(tc.Ctx, ac)

	if err != nil {
		logger.Error("Error rendering state", zap.Error(err))
	}

}

func DefaultHandleMessage[S any, C any](ac *ApplicationChat[S, C], tc *telegram.TelegramContextTextMessage) {
	logger := ac.Loggers.Handle.With(zap.Int64("UpdateID", tc.UpdateID))

	logger.Info("HandleMessage", zap.Any("text", tc.Text))
	// tc.Logger.Debug("LocalStateTree", zap.String("tree", ac.State.treeState.LocalStateTree.String()))

	ac.State.LockState(logger.Named("LockState"))
	defer ac.State.UnlockState(logger.Named("LockState"))

	if ac.State.inputHandler != nil {

		// tc.Logger.Debug("HandleMessage", zap.Any("message", tc.Message))

		ac.State.renderedElements = append(
			ac.State.renderedElements,
			render.NewRenderedUserMessage(tc.Message.ID),
		)

		action := ac.State.inputHandler(tc.Message.Message)

		internalActionHandle(ac, &tc.TelegramUpdateContext, action)

	} else {
		logger.Warn("Missing InputHandler")
	}

	err := ac.App.RenderFunc(tc.Ctx, ac)

	if err != nil {
		logger.Error("Error rendering state", zap.Error(err))
	}
}

// Handle
func DefaultHandleActionExternal[S any, C any](ac *ApplicationChat[S, C], tc *telegram.TelegramUpdateContext, action any) {

	actionName := reflection.ReflectStructName(action)
	logger := ac.Loggers.Action.With(zap.String("action", actionName))

	logger.Info("HandleActionExternal")

	ac.State.LockState(logger.Named("LockState"))
	defer ac.State.UnlockState(logger.Named("LockState"))

	internalActionHandle(ac, tc, action)

	err := ac.App.RenderFunc(tc.Ctx, ac)

	if err != nil {
		logger.Error("Error rendering state", zap.Error(err))
	}

}
