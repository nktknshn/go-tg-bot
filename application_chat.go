package tgbot

import (
	"sync"

	"go.uber.org/zap"
)

// tie together an application methods and a chat state
type ApplicationChat[S any, C any] struct {
	App   *Application[S, C]
	State *ChatState[S, C]

	// logger with attached chatID
	Logger *zap.Logger
}

func NewApplicationChat[S any, C any](app Application[S, C], tc *TelegramContext) *ApplicationChat[S, C] {
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

	res := app.PreRender(&chatState)

	return &ApplicationChat[S, C]{
		App:   &app,
		State: &res.InternalChatState,
		Logger: app.Loggers.ApplicationChat(
			GetLogger().With(zap.Int64("ChatID", tc.ChatID)),
		),
	}
}

// Computes the output based on the state and renders it to the user
func DefaultRenderFunc[S any, C any](ac *ApplicationChat[S, C]) error {
	ac.Logger.Info("RenderFunc")

	res := ac.App.PreRender(ac.State)
	rendered, err := res.ExecuteRender(ac.State.Renderer)

	if err != nil {
		ac.Logger.Error("Error in RenderFunc", zap.Error(err))
		return err
	}

	ac.State = &res.InternalChatState
	ac.State.renderedElements = rendered

	return nil
}

func DefaultHandlerCallback[S any, C any](ac *ApplicationChat[S, C], tc *TelegramContextCallback) {
	tc.Logger.Info("HandleCallback", zap.Any("data", tc.UpdateBotCallbackQuery.QueryID))
	tc.Logger.Debug("LocalStateTree", zap.String("tree", ac.State.treeState.LocalStateTree.String()))

	ac.State.LockState(tc.Logger)
	defer ac.State.UnlockState(tc.Logger)

	if ac.State.callbackHandler != nil {
		result := ac.State.callbackHandler(string(tc.UpdateBotCallbackQuery.Data))

		ac.Logger.Debug("HandleCallback", zap.Any("action", result))

		if result == nil {
			return
		}

		internalActionHandle(ac, &tc.TelegramContext, result.action)

		if !result.noCallback {
			tc.AnswerCallbackQuery()
		}

	} else {
		tc.Logger.Warn("Missing CallbackHandler")
	}

	err := ac.App.RenderFunc(ac)

	if err != nil {
		tc.Logger.Error("Error rendering state", zap.Error(err))
	}

}

func DefaultHandleMessage[S any, C any](ac *ApplicationChat[S, C], tc *TelegramContextTextMessage) {

	tc.Logger.Info("HandleMessage", zap.Any("text", tc.Text))
	tc.Logger.Debug("LocalStateTree", zap.String("tree", ac.State.treeState.LocalStateTree.String()))

	ac.State.LockState(tc.Logger)
	defer ac.State.UnlockState(tc.Logger)

	if ac.State.inputHandler != nil {

		tc.Logger.Debug("HandleMessage", zap.Any("message", tc.Message))

		ac.State.renderedElements = append(
			ac.State.renderedElements,
			newRenderedUserMessage(tc.Message.ID),
		)

		action := ac.State.inputHandler(tc.Message.Message)

		internalActionHandle(ac, &tc.TelegramContext, action)

	} else {
		tc.Logger.Warn("Missing InputHandler")
	}

	err := ac.App.RenderFunc(ac)

	if err != nil {
		tc.Logger.Error("Error rendering state", zap.Error(err))
	}
}

func DefaultHandleActionExternal[S any, C any](ac *ApplicationChat[S, C], tc *TelegramContext, action any) {
	ac.Logger.Info("HandleActionExternal", zap.String("action", reflectStructName(action)))

	ac.State.LockState(tc.Logger)
	defer ac.State.UnlockState(tc.Logger)

	internalActionHandle(ac, tc, action)

	err := ac.App.RenderFunc(ac)

	if err != nil {
		tc.Logger.Error("Error rendering state", zap.Error(err))
	}

}
