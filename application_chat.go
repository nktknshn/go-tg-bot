package tgbot

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

type ApplicationChatLoggers struct {
	Root      *zap.Logger
	Component *zap.Logger
	Update    *zap.Logger
	Action    *zap.Logger
	Render    *zap.Logger
}

// tie together an application methods and a chat state
type ApplicationChat[S any, C any] struct {
	App   *Application[S, C]
	State *ChatState[S, C]

	// loggers for different parts of the app
	Loggers *ApplicationChatLoggers
}

type ExecutionContext struct {
	UpdateContext *TelegramUpdateContext
}

func NewApplicationChat[S any, C any](app Application[S, C], tc *TelegramUpdateContext) *ApplicationChat[S, C] {
	appState := app.CreateAppState(tc)

	chatState := ChatState[S, C]{
		ChatID:           tc.ChatID,
		AppState:         appState,
		renderedElements: []RenderedElement{},
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
		Update: rootLogger.Named("Update"),
		Action: rootLogger.Named("Action"),
		Render: rootLogger.Named("Render"),
	}

	res := app.ComputeNextState(&chatState, loggers.Component)

	return &ApplicationChat[S, C]{
		App:     &app,
		State:   &res.NextChatState,
		Loggers: loggers,
	}
}

// Computes the output based on the state and renders it to the user
func DefaultRenderFunc[S any, C any](ctx context.Context, ac *ApplicationChat[S, C]) error {
	ac.Loggers.Render.Debug("RenderFunc called")

	res := ac.App.ComputeNextState(ac.State, ac.Loggers.Component)
	rendered, err := ExecuteRenderActions(ctx, ac.State.Renderer, res.RenderActions, ac.Loggers.Render)

	if err != nil {
		ac.Loggers.Root.Error("Error in RenderFunc", zap.Error(err))
		return err
	}

	ac.State = &res.NextChatState
	ac.State.renderedElements = rendered

	return nil
}

func DefaultHandlerCallback[S any, C any](ac *ApplicationChat[S, C], tc *TelegramContextCallback) {

	logger := ac.Loggers.Update.With(zap.Int64("UpdateID", tc.UpdateID))

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

		internalActionHandle(ac, &tc.TelegramUpdateContext, result.action, logger.Named("Action"))

		if !result.noAnswer {
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

func DefaultHandleMessage[S any, C any](ac *ApplicationChat[S, C], tc *TelegramContextTextMessage) {
	logger := ac.Loggers.Update.With(zap.Int64("UpdateID", tc.UpdateID))

	logger.Info("HandleMessage", zap.Any("text", tc.Text))
	// tc.Logger.Debug("LocalStateTree", zap.String("tree", ac.State.treeState.LocalStateTree.String()))

	ac.State.LockState(logger.Named("LockState"))
	defer ac.State.UnlockState(logger.Named("LockState"))

	if ac.State.inputHandler != nil {

		// tc.Logger.Debug("HandleMessage", zap.Any("message", tc.Message))

		ac.State.renderedElements = append(
			ac.State.renderedElements,
			newRenderedUserMessage(tc.Message.ID),
		)

		action := ac.State.inputHandler(tc.Message.Message)

		internalActionHandle(ac, &tc.TelegramUpdateContext, action, logger)

	} else {
		logger.Warn("Missing InputHandler")
	}

	err := ac.App.RenderFunc(tc.Ctx, ac)

	if err != nil {
		logger.Error("Error rendering state", zap.Error(err))
	}
}

// Handle
func DefaultHandleActionExternal[S any, C any](ac *ApplicationChat[S, C], tc *TelegramUpdateContext, action any) {

	actionName := reflectStructName(action)
	logger := ac.Loggers.Action.With(zap.String("action", actionName))

	logger.Info("HandleActionExternal")

	ac.State.LockState(logger.Named("LockState"))
	defer ac.State.UnlockState(logger.Named("LockState"))

	internalActionHandle(ac, tc, action, logger)

	err := ac.App.RenderFunc(tc.Ctx, ac)

	if err != nil {
		logger.Error("Error rendering state", zap.Error(err))
	}

}
