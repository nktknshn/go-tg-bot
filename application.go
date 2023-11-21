package tgbot

import (
	"go.uber.org/zap"
)

type callbackResult struct {
	action any

	//
	noCallback bool
}

// Handles text input from the user
// Returnes action that is going to be dispatched to the application
// Returns Next if no action is needed
type chatInputHandler func(string) any

// Handle inline button click. Returns nil if no action has matched the callback data
type chatCallbackHandler func(string) *callbackResult

// User defined function
type handleMessageFunc[S any, C any] func(*ApplicationContext[S, C], *TelegramContext)
type handleCallbackFunc[S any, C any] func(*ApplicationContext[S, C], *TelegramContext)

type handleInitFunc[S any] func(*TelegramContext)

type handleActionFunc[S any, C any] func(*ApplicationContext[S, C], *TelegramContext, any)

type renderFuncType[S any, C any] func(*ApplicationContext[S, C]) error

type stateToCompFuncType[S any, C any] func(S) Comp

type ApplicationContext[S any, C any] struct {
	App    *Application[S, C]
	State  *ChatState[S, C]
	Logger *zap.Logger
}

// Defines Application with state S
type Application[S any, C any] struct {
	CreateAppState func(*TelegramContext) S

	HandleActionExternal handleActionFunc[S, C]

	// actions reducer
	HandleAction handleActionFunc[S, C]

	HandleMessage handleMessageFunc[S, C]

	HandleCallback handleCallbackFunc[S, C]

	CreateGlobalContext func(state *ChatState[S, C]) C

	HandleInit handleInitFunc[S]

	// use state to render bot interface to the user
	RenderFunc  renderFuncType[S, C]
	StateToComp stateToCompFuncType[S, C]

	CreateChatRenderer func(*TelegramContext) ChatRenderer
}

type NewApplicationProps[S any, C any] struct {
	HandleMessage       handleMessageFunc[S, C]
	HandleCallback      handleCallbackFunc[S, C]
	HandleInit          handleInitFunc[S]
	RenderFunc          renderFuncType[S, C]
	CreateRenderer      func(*TelegramContext) ChatRenderer
	CreateGlobalContext func(*ChatState[S, C]) C
}

// type EmptyGlobalContext struct{}

// func (egc EmptyGlobalContext) Get() any {
// 	panic("EmptyGlobalContext")
// }

// func DefaultCreateContext[S any, C any](state *ChatState[S, C]) C {
// 	return nil
// }

func DefaultHandleActionExternal[S any, C any](ac *ApplicationContext[S, C], tc *TelegramContext, action any) {
	ac.Logger.Info("HandleActionExternal", zap.String("action", reflectStructName(action)))

	ac.State.LockState(tc.Logger)
	defer ac.State.UnlockState(tc.Logger)

	internalHandleAction(ac, tc, action)

	err := ac.App.RenderFunc(ac)

	if err != nil {
		tc.Logger.Error("Error rendering state", zap.Error(err))
	}

}

// Computes the output based on the state and renders it to the user
func DefaultRenderFunc[S any, C any](ac *ApplicationContext[S, C]) error {
	ac.Logger.Info("RenderFunc")

	res := ac.App.PreRender(ac)
	rendered, err := res.ExecuteRender(ac.State.Renderer)

	if err != nil {
		ac.Logger.Error("Error in RenderFunc", zap.Error(err))
		return err
	}

	ac.State = &res.InternalChatState
	ac.State.renderedElements = rendered

	return nil
}

func DefaultHandlerCallback[S any, C any](ac *ApplicationContext[S, C], tc *TelegramContext) {
	tc.Logger.Info("HandleCallback", zap.Any("data", tc.Update.CallbackQuery.Data))
	tc.Logger.Debug("LocalStateTree", zap.String("tree", ac.State.treeState.LocalStateTree.String()))

	ac.State.LockState(tc.Logger)
	defer ac.State.UnlockState(tc.Logger)

	if ac.State.callbackHandler != nil {
		result := ac.State.callbackHandler(tc.Update.CallbackQuery.Data)

		ac.Logger.Debug("HandleCallback", zap.Any("action", result))

		if result == nil {
			return
		}

		internalHandleAction(ac, tc, result.action)

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

func DefaultHandleMessage[S any, C any](ac *ApplicationContext[S, C], tc *TelegramContext) {

	tc.Logger.Info("HandleMessage", zap.Any("text", tc.Update.Message.Text))
	tc.Logger.Debug("LocalStateTree", zap.String("tree", ac.State.treeState.LocalStateTree.String()))

	ac.State.LockState(tc.Logger)
	defer ac.State.UnlockState(tc.Logger)

	if ac.State.inputHandler != nil {

		ac.State.renderedElements = append(
			ac.State.renderedElements, newRenderedUserMessage(tc.Update.Message.ID),
		)

		action := ac.State.inputHandler(tc.Update.Message.Text)

		internalHandleAction(ac, tc, action)

	} else {
		tc.Logger.Warn("Missing InputHandler")
	}

	err := ac.App.RenderFunc(ac)

	if err != nil {
		tc.Logger.Error("Error rendering state", zap.Error(err))
	}
}

func NewApplication[S any, C any](
	// Creates state
	createAppState func(*TelegramContext) S,
	// turns state into basic elements
	stateToComp stateToCompFuncType[S, C],
	// handles action
	handleAction handleActionFunc[S, C],
	propss ...*NewApplicationProps[S, C],
) *Application[S, C] {

	props := &NewApplicationProps[S, C]{}

	if len(propss) > 0 {
		props = propss[0]
	}

	var (
		handleMessage  = props.HandleMessage
		handleCallback = props.HandleCallback
		handleInit     = props.HandleInit
		renderFunc     = props.RenderFunc
		createRenderer = props.CreateRenderer
		createContext  = props.CreateGlobalContext
	)

	if handleMessage == nil {
		handleMessage = DefaultHandleMessage[S, C]
	}

	if handleCallback == nil {
		handleCallback = DefaultHandlerCallback[S, C]
	}

	if handleInit == nil {
		handleInit = func(tc *TelegramContext) {}
	}

	if renderFunc == nil {
		renderFunc = DefaultRenderFunc[S, C]
	}

	if createRenderer == nil {
		createRenderer = func(tc *TelegramContext) ChatRenderer {
			return NewTelegramChatRenderer(tc.Bot, tc.ChatID)
		}
	}

	// if createContext == nil {
	// 	createContext = DefaultCreateContext[S, C]
	// }

	return &Application[S, C]{
		CreateAppState:       createAppState,
		StateToComp:          stateToComp,
		HandleAction:         handleAction,
		HandleMessage:        handleMessage,
		HandleCallback:       handleCallback,
		HandleInit:           handleInit,
		RenderFunc:           renderFunc,
		CreateChatRenderer:   createRenderer,
		CreateGlobalContext:  createContext,
		HandleActionExternal: DefaultHandleActionExternal[S, C],
	}
}
