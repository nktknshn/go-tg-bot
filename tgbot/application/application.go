package application

import (
	"context"

	"github.com/nktknshn/go-tg-bot/tgbot/component"
	"github.com/nktknshn/go-tg-bot/tgbot/logging"
	"github.com/nktknshn/go-tg-bot/tgbot/render"
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
	"go.uber.org/zap"
)

// Handles text input from the user
// Returnes action that is going to be dispatched to the application

// Handle inline button click. Returns nil if no action has matched the callback data

// User defined function
type HandleMessageFunc[S any, C any] func(*ApplicationChat[S, C], *telegram.TelegramContextTextMessage)
type HandleCallbackFunc[S any, C any] func(*ApplicationChat[S, C], *telegram.TelegramContextCallback)

type HandleInitFunc[S any] func(*telegram.TelegramUpdateContext)

type HandleActionFunc[S any, C any] func(*ApplicationChat[S, C], *telegram.TelegramUpdateContext, any)

type RenderFuncType[S any, C any] func(context.Context, *ApplicationChat[S, C]) error

type StateToCompFuncType[S any, C any] func(S) component.Comp

type CreateAppStateFunc[S, C any] func(*Application[S, C], *telegram.TelegramUpdateContext, *zap.Logger) S

// Defines Application with state S
type Application[S any, C any] struct {
	CreateAppState CreateAppStateFunc[S, C]

	HandleActionExternal HandleActionFunc[S, C]

	// actions reducer
	HandleAction HandleActionFunc[S, C]

	HandleMessage HandleMessageFunc[S, C]

	HandleCallback HandleCallbackFunc[S, C]

	CreateGlobalContext func(state *ApplicationChatState[S, C]) C

	// not used currently
	HandleInit HandleInitFunc[S]

	StateToComp StateToCompFuncType[S, C]

	// use state to render bot interface to the user
	RenderFunc RenderFuncType[S, C]

	CreateChatRenderer func(*telegram.TelegramUpdateContext) render.ChatRenderer

	Loggers logging.Loggers
}

type NewApplicationProps[S any, C any] struct {
	HandleMessage       HandleMessageFunc[S, C]
	HandleCallback      HandleCallbackFunc[S, C]
	HandleInit          HandleInitFunc[S]
	RenderFunc          RenderFuncType[S, C]
	CreateRenderer      func(*telegram.TelegramUpdateContext) render.ChatRenderer
	CreateGlobalContext func(*ApplicationChatState[S, C]) C
}

func (app *Application[S, C]) globalContext(chatState *ApplicationChatState[S, C]) component.GlobalContext[C] {
	if app.CreateGlobalContext != nil {
		ctxValue := app.CreateGlobalContext(chatState)
		return component.NewGlobalContextTyped[C](ctxValue)
	} else {
		return component.NewEmptyGlobalContext()
	}
}

func New[S any, C any](
	// Creates state
	createAppState func(*Application[S, C], *telegram.TelegramUpdateContext, *zap.Logger) S,
	// turns state into basic elements
	stateToComp StateToCompFuncType[S, C],
	// handles action
	handleAction HandleActionFunc[S, C],
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
		handleInit = func(tc *telegram.TelegramUpdateContext) {}
	}

	if renderFunc == nil {
		renderFunc = DefaultRenderFunc[S, C]
	}

	if createRenderer == nil {
		createRenderer = func(tc *telegram.TelegramUpdateContext) render.ChatRenderer {
			return render.NewTelegramChatRenderer(tc.Bot, tc.Update.User)
		}
	}

	loggers := &logging.LoggersDefault{}

	return &Application[S, C]{
		HandleInit:           handleInit,
		CreateAppState:       createAppState,
		StateToComp:          stateToComp,
		HandleAction:         handleAction,
		HandleMessage:        handleMessage,
		HandleCallback:       handleCallback,
		RenderFunc:           renderFunc,
		CreateChatRenderer:   createRenderer,
		CreateGlobalContext:  createContext,
		HandleActionExternal: DefaultHandleActionExternal[S, C],
		Loggers:              loggers,
	}
}
