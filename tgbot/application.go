package tgbot

import (
	"context"

	"github.com/nktknshn/go-tg-bot/tgbot/component"
	"github.com/nktknshn/go-tg-bot/tgbot/logging"
)

// Handles text input from the user
// Returnes action that is going to be dispatched to the application

// Handle inline button click. Returns nil if no action has matched the callback data

// User defined function
type handleMessageFunc[S any, C any] func(*ApplicationChat[S, C], *TelegramContextTextMessage)
type handleCallbackFunc[S any, C any] func(*ApplicationChat[S, C], *TelegramContextCallback)

type handleInitFunc[S any] func(*TelegramUpdateContext)

type handleActionFunc[S any, C any] func(*ApplicationChat[S, C], *TelegramUpdateContext, any)

type renderFuncType[S any, C any] func(context.Context, *ApplicationChat[S, C]) error

type stateToCompFuncType[S any, C any] func(S) component.Comp

type createAppStateFunc[S any] func(*TelegramUpdateContext) S

// Defines Application with state S
type Application[S any, C any] struct {
	CreateAppState createAppStateFunc[S]

	HandleActionExternal handleActionFunc[S, C]

	// actions reducer
	HandleAction handleActionFunc[S, C]

	HandleMessage handleMessageFunc[S, C]

	HandleCallback handleCallbackFunc[S, C]

	CreateGlobalContext func(state *ChatState[S, C]) C

	// not used currently
	HandleInit handleInitFunc[S]

	StateToComp stateToCompFuncType[S, C]

	// use state to render bot interface to the user
	RenderFunc renderFuncType[S, C]

	CreateChatRenderer func(*TelegramUpdateContext) ChatRenderer

	Loggers logging.TgbotLoggers
}

type ApplicationProps[S any, C any] struct {
	HandleMessage       handleMessageFunc[S, C]
	HandleCallback      handleCallbackFunc[S, C]
	HandleInit          handleInitFunc[S]
	RenderFunc          renderFuncType[S, C]
	CreateRenderer      func(*TelegramUpdateContext) ChatRenderer
	CreateGlobalContext func(*ChatState[S, C]) C
}

func NewApplication[S any, C any](
	// Creates state
	createAppState func(*TelegramUpdateContext) S,
	// turns state into basic elements
	stateToComp stateToCompFuncType[S, C],
	// handles action
	handleAction handleActionFunc[S, C],
	propss ...*ApplicationProps[S, C],
) *Application[S, C] {

	props := &ApplicationProps[S, C]{}

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
		handleInit = func(tc *TelegramUpdateContext) {}
	}

	if renderFunc == nil {
		renderFunc = DefaultRenderFunc[S, C]
	}

	if createRenderer == nil {
		createRenderer = func(tc *TelegramUpdateContext) ChatRenderer {
			return NewTelegramChatRenderer(tc.Bot, tc.Update.User)
		}
	}

	loggers := logging.DefaultLoggers

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

func (app *Application[S, C]) globalContext(chatState *ChatState[S, C]) component.GlobalContext[C] {
	if app.CreateGlobalContext != nil {
		ctxValue := app.CreateGlobalContext(chatState)
		return component.NewGlobalContextTyped[C](ctxValue)
	} else {
		return component.NewEmptyGlobalContext()
	}
}

func (a *Application[S, C]) NewHandler(tc *TelegramUpdateContext) *ChatHandlerImpl[S, C] {
	return NewChatHandler[S, C](*a, tc)
}

func (a *Application[S, C]) ChatsDispatcher() *ChatsDispatcher {

	return NewChatsDispatcher(&ChatsDispatcherProps{
		ChatFactory: &factoryFunc{
			f: func(tc *TelegramUpdateContext) ChatHandler {
				return a.NewHandler(tc)
			},
		},
	})
}

type factoryFunc struct {
	f func(*TelegramUpdateContext) ChatHandler
}

func (f *factoryFunc) CreateChatHandler(tc *TelegramUpdateContext) ChatHandler {
	return f.f(tc)
}
