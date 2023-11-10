package tgbot

import (
	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

// A is a type of returned Action to be used in actions reducers
type ChatInputHandler[A any] func(string) A
type ChatCallbackHandler[A any] func(string) *A

type InternalChatState[S any, A any] struct {
	ChatID int64
	// state of the application
	AppState S

	// elements visible to the user
	RenderedElements []RenderedElement

	// handler for text messages
	InputHandler ChatInputHandler[A]

	// handler for callback queries
	CallbackHandler ChatCallbackHandler[A]

	Renderer ChatRenderer
}

type HandleMessageFunc[S any, A any] func(*ApplicationContext[S, A], *TelegramContext)
type HandleCallbackFunc[S any, A any] func(*ApplicationContext[S, A], *TelegramContext)
type HandleInitFunc[S any] func(*TelegramContext)

type HandleActionFunc[S any, A any] func(*ApplicationContext[S, A], *TelegramContext, A)

type RenderFuncType[S any, A any] func(*ApplicationContext[S, A]) error
type StateToCompFuncType[S any, A any] func(S) Comp[A]

type ApplicationContext[S any, A any] struct {
	App    *Application[S, A]
	State  *InternalChatState[S, A]
	Logger *zap.Logger
}

// Defines Application with state S
type Application[S any, A any] struct {
	CreateAppState func(*TelegramContext) S

	// actions reducer
	HandleAction HandleActionFunc[S, A]

	HandleMessage HandleMessageFunc[S, A]

	HandleCallback HandleCallbackFunc[S, A]

	// HandleEvent

	HandleInit HandleInitFunc[S]

	// taken S renderes elements
	RenderFunc  RenderFuncType[S, A]
	StateToComp StateToCompFuncType[S, A]

	CreateChatRenderer func(*TelegramContext) ChatRenderer
}

type NewApplicationProps[S any, A any] struct {
	// CreateAppState func(*TelegramContext) S
	// HandleAction  HandleActionFunc[S, A]
	HandleMessage  HandleMessageFunc[S, A]
	HandleCallback HandleCallbackFunc[S, A]
	HandleInit     HandleInitFunc[S]
	RenderFunc     RenderFuncType[S, A]
	CreateRenderer func(*TelegramContext) ChatRenderer
}

func DefaultRenderFunc[S any, A any](ac *ApplicationContext[S, A]) error {
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

func handleAction[S any, A any](ac *ApplicationContext[S, A], tc *TelegramContext, a A) {
	tc.Logger.Debug("HandleAction", zap.Any("action", a))
	ac.App.HandleAction(ac, tc, a)
}

func DefaultHandlerCallback[S any, A any](ac *ApplicationContext[S, A], tc *TelegramContext) {
	tc.Logger.Info("HandleCallback", zap.Any("data", tc.Update.CallbackQuery.Data))

	if ac.State.CallbackHandler != nil {
		action := ac.State.CallbackHandler(tc.Update.CallbackQuery.Data)

		ac.Logger.Debug("HandleCallback", zap.Any("action", action))

		if action == nil {
			return
		}

		handleAction(ac, tc, *action)

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

func DefaultHandleMessage[S any, A any](ac *ApplicationContext[S, A], tc *TelegramContext) {
	tc.Logger.Info("HandleMessage", zap.Any("text", tc.Update.Message.Text))

	if ac.State.InputHandler != nil {

		ac.State.RenderedElements = append(ac.State.RenderedElements, NewRenderedUserMessage(tc.Update.Message.ID))

		action := ac.State.InputHandler(tc.Update.Message.Text)

		handleAction(ac, tc, action)

		err := ac.App.RenderFunc(ac)

		if err != nil {
			tc.Logger.Error("Error rendering state", zap.Error(err))
		}

	} else {
		tc.Logger.Error("Missing InputHandler")
	}
}

func NewApplication[S any, A any](
	// Creates state
	createAppState func(*TelegramContext) S,
	// handles action
	handleAction HandleActionFunc[S, A],
	// turns state into basic elements
	stateToComp StateToCompFuncType[S, A],
	propss ...*NewApplicationProps[S, A],
) *Application[S, A] {

	props := &NewApplicationProps[S, A]{}

	if len(propss) > 0 {
		props = propss[0]
	}

	var (
		handleMessage  = props.HandleMessage
		handleCallback = props.HandleCallback
		handleInit     = props.HandleInit
		renderFunc     = props.RenderFunc
		createRenderer = props.CreateRenderer
	)

	if handleMessage == nil {
		handleMessage = DefaultHandleMessage[S, A]
	}

	if handleCallback == nil {
		handleCallback = DefaultHandlerCallback[S, A]
	}

	if handleInit == nil {
		handleInit = func(tc *TelegramContext) {}
	}

	if renderFunc == nil {
		renderFunc = DefaultRenderFunc[S, A]
	}

	if createRenderer == nil {
		createRenderer = func(tc *TelegramContext) ChatRenderer {
			return NewTelegramChatRenderer(tc.Bot, tc.ChatID)
		}
	}

	return &Application[S, A]{
		CreateAppState:     createAppState,
		StateToComp:        stateToComp,
		HandleAction:       handleAction,
		HandleMessage:      handleMessage,
		HandleCallback:     handleCallback,
		HandleInit:         handleInit,
		RenderFunc:         renderFunc,
		CreateChatRenderer: createRenderer,
	}
}

func (a *Application[S, A]) NewHandler(tc *TelegramContext) *Handler[S, A] {
	return NewHandler[S, A](*a, tc)
}
