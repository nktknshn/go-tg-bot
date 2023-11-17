package tgbot

import (
	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

// A is a type of returned Action to be used in actions reducers
type ChatInputHandler[A any] func(string) A
type ChatCallbackHandler[A any] func(string) *A

type InternalChatState[S any, A any, C any] struct {
	ChatID int64

	// state of the application
	AppState S

	// state of the application
	TreeState *RunResultWithStateTree[A]

	// elements visible to the user
	RenderedElements []RenderedElement

	// handler for text messages
	InputHandler ChatInputHandler[any]

	// handler for callback queries
	CallbackHandler ChatCallbackHandler[A]

	Renderer ChatRenderer
}

type HandleMessageFunc[S any, A any, C any] func(*ApplicationContext[S, A, C], *TelegramContext)
type HandleCallbackFunc[S any, A any, C any] func(*ApplicationContext[S, A, C], *TelegramContext)
type HandleInitFunc[S any] func(*TelegramContext)

type HandleActionFunc[S any, A any, C any] func(*ApplicationContext[S, A, C], *TelegramContext, any)

type RenderFuncType[S any, A any, C any] func(*ApplicationContext[S, A, C]) error
type StateToCompFuncType[S any, A any, C any] func(S) Comp[A]

type ApplicationContext[S any, A any, C any] struct {
	App    *Application[S, A, C]
	State  *InternalChatState[S, A, C]
	Logger *zap.Logger
}

// Defines Application with state S
type Application[S any, A any, C any] struct {
	CreateAppState func(*TelegramContext) S

	// actions reducer
	HandleAction HandleActionFunc[S, A, C]

	HandleMessage HandleMessageFunc[S, A, C]

	HandleCallback HandleCallbackFunc[S, A, C]

	CreateGlobalContext func(state *InternalChatState[S, A, C]) GlobalContextTyped[C]

	HandleInit HandleInitFunc[S]

	// taken S renderes elements
	RenderFunc  RenderFuncType[S, A, C]
	StateToComp StateToCompFuncType[S, A, C]

	CreateChatRenderer func(*TelegramContext) ChatRenderer
}

type NewApplicationProps[S any, A any, C any] struct {
	// CreateAppState func(*TelegramContext) S
	// HandleAction  HandleActionFunc[S, A, C]
	HandleMessage       HandleMessageFunc[S, A, C]
	HandleCallback      HandleCallbackFunc[S, A, C]
	HandleInit          HandleInitFunc[S]
	RenderFunc          RenderFuncType[S, A, C]
	CreateRenderer      func(*TelegramContext) ChatRenderer
	CreateGlobalContext func(*InternalChatState[S, A, C]) GlobalContextTyped[C]
}

func DefaultCreateContext[S any, A any, C any](state *InternalChatState[S, A, C]) GlobalContextTyped[C] {
	return nil
}

func DefaultRenderFunc[S any, A any, C any](ac *ApplicationContext[S, A, C]) error {
	ac.Logger.Info("RenderFunc")

	res := ac.App.PreRender(ac)
	rendered, err := res.ExecuteRender(ac.State.Renderer)

	if err != nil {
		ac.Logger.Error("Error in RenderFunc", zap.Error(err))
		return err
	}

	ac.State = &res.InternalChatState
	ac.State.RenderedElements = rendered

	return nil
}

func DefaultHandlerCallback[S any, A any, C any](ac *ApplicationContext[S, A, C], tc *TelegramContext) {
	tc.Logger.Info("HandleCallback", zap.Any("data", tc.Update.CallbackQuery.Data))
	tc.Logger.Debug("LocalStateTree", zap.String("tree", ac.State.TreeState.LocalStateTree.String()))

	if ac.State.CallbackHandler != nil {
		action := ac.State.CallbackHandler(tc.Update.CallbackQuery.Data)

		ac.Logger.Debug("HandleCallback", zap.Any("action", action))

		if action == nil {
			return
		}

		internalHandleAction(ac, tc, *action)

		tc.Bot.AnswerCallbackQuery(tc.Ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: tc.Update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		err := ac.App.RenderFunc(ac)

		if err != nil {
			tc.Logger.Error("Error rendering state", zap.Error(err))
		}
	} else {
		tc.Logger.Error("Missing CallbackHandler")
	}

}

func DefaultHandleMessage[S any, A any, C any](ac *ApplicationContext[S, A, C], tc *TelegramContext) {
	tc.Logger.Info("HandleMessage", zap.Any("text", tc.Update.Message.Text))
	tc.Logger.Debug("LocalStateTree", zap.String("tree", ac.State.TreeState.LocalStateTree.String()))

	if ac.State.InputHandler != nil {

		ac.State.RenderedElements = append(
			ac.State.RenderedElements, NewRenderedUserMessage(tc.Update.Message.ID),
		)

		action := ac.State.InputHandler(tc.Update.Message.Text)

		internalHandleAction(ac, tc, action)

		err := ac.App.RenderFunc(ac)

		if err != nil {
			tc.Logger.Error("Error rendering state", zap.Error(err))
		}

	} else {
		tc.Logger.Error("Missing InputHandler")
	}
}

func NewApplication[S any, A any, C any](
	// Creates state
	createAppState func(*TelegramContext) S,
	// turns state into basic elements
	stateToComp StateToCompFuncType[S, A, C],
	// handles action
	handleAction HandleActionFunc[S, A, C],
	propss ...*NewApplicationProps[S, A, C],
) *Application[S, A, C] {

	props := &NewApplicationProps[S, A, C]{}

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
		handleMessage = DefaultHandleMessage[S, A, C]
	}

	if handleCallback == nil {
		handleCallback = DefaultHandlerCallback[S, A, C]
	}

	if handleInit == nil {
		handleInit = func(tc *TelegramContext) {}
	}

	if renderFunc == nil {
		renderFunc = DefaultRenderFunc[S, A, C]
	}

	if createRenderer == nil {
		createRenderer = func(tc *TelegramContext) ChatRenderer {
			return NewTelegramChatRenderer(tc.Bot, tc.ChatID)
		}
	}

	if createContext == nil {
		createContext = DefaultCreateContext[S, A, C]
	}

	return &Application[S, A, C]{
		CreateAppState:      createAppState,
		StateToComp:         stateToComp,
		HandleAction:        handleAction,
		HandleMessage:       handleMessage,
		HandleCallback:      handleCallback,
		HandleInit:          handleInit,
		RenderFunc:          renderFunc,
		CreateChatRenderer:  createRenderer,
		CreateGlobalContext: createContext,
	}
}

func (a *Application[S, A, C]) NewHandler(tc *TelegramContext) *Handler[S, A, C] {
	return NewHandler[S, A, C](*a, tc)
}

func (a *Application[S, A, C]) ChatsDispatcher() *ChatsDispatcher {
	return NewChatsDispatcher(&ChatsDispatcherProps{
		ChatFactory: func(tc *TelegramContext) ChatHandler {
			return a.NewHandler(tc)
		},
	})
}
