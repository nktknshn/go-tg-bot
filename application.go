package tgbot

import (
	"context"

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

// func (s *InternalChatState[S]) ModifyState() {
// 	return
// }

// type renderFunc[S any] func(S) *renderFuncResult[S]

// type renderFuncResult[S any] struct {
// 	chatState *S
// }

type HandleMessageFunc[S any, A any] func(*ApplicationContext[S, A], *TelegramContext)
type HandleCallbackFunc[S any, A any] func(*ApplicationContext[S, A], *TelegramContext)
type HandleInitFunc[S any] func(*TelegramContext)

type HandleActionFunc[S any, A any] func(*ApplicationContext[S, A], *TelegramContext, A)

// type ReducerFuncType[A any, S any] func(InternalChatState[S]) InternalChatState[S]

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

func (a *Application[S, A]) NewHandler(tc *TelegramContext) *Handler[S, A] {
	return NewHandler[S, A](*a, tc)
}

type PreRenderData[S any, A any] struct {
	InternalChatState InternalChatState[S, A]
	ExecuteRender     func(renderer ChatRenderer) ([]RenderedElement, error)
}

func getCallbackHandlersMap[A any](outcomingMessages []OutcomingMessage) map[string]func() *A {

	callbackHandlers := make(map[string]func() *A)

	for _, m := range outcomingMessages {
		switch el := m.(type) {
		case *OutcomingTextMessage[A]:
			for _, row := range el.Buttons {
				for _, butt := range row {

					logger.Info("Setting callback handler", zap.String("key", butt.CallbackData()))
					butt := butt
					callbackHandlers[butt.CallbackData()] = func() *A {
						v := butt.OnClick()
						return &v
					}
				}
			}
		}
	}

	return callbackHandlers
}

func callbackMapToHandler[A any](cbmap map[string]func() *A) ChatCallbackHandler[A] {
	return func(callbackData string) *A {

		logger.Info("Callback handler", zap.String("data", callbackData))

		if handler, ok := cbmap[callbackData]; ok {
			// ac.App.HandleAction(ac, tc, )
			logger.Info("Calling handler", zap.String("data", callbackData))

			return handler()
		} else {
			// err := fmt.Errorf("no handler for callback %v", key)
			logger.Error("No handler for callback", zap.String("key", callbackData))
			return nil
		}

	}
}

func (a *Application[S, A]) PreRender(ac *ApplicationContext[S, A]) *PreRenderData[S, A] {

	els := ComponentToElements2(
		ac.App.StateToComp(ac.State.AppState),
	)

	res := ElementsToMessagesAndHandlers[A](els)

	ac.Logger.Debug("PreRender", zap.Any("res", res))

	var inputHandler ChatInputHandler[A]
	callbackMap := getCallbackHandlersMap[A](res.OutcomingMessages)
	callbackHandler := callbackMapToHandler[A](callbackMap)

	if len(res.InputHandlers) > 0 {
		inputHandler = func(text string) A {
			return res.InputHandlers[0].Handler(text)
		}
	}

	nextState := InternalChatState[S, A]{
		ChatID:           ac.State.ChatID,
		AppState:         ac.State.AppState,
		RenderedElements: ac.State.RenderedElements,
		InputHandler:     inputHandler,
		CallbackHandler:  callbackHandler,
		Renderer:         ac.State.Renderer,
	}

	return &PreRenderData[S, A]{
		InternalChatState: nextState,
		ExecuteRender: func(renderer ChatRenderer) ([]RenderedElement, error) {
			ac.Logger.Info("ExecuteRender")

			actions := GetRenderActions[A](
				ac.State.RenderedElements,
				res.OutcomingMessages,
			)

			// logger.Info("RenderActions", zap.Any("actions", actions))

			for _, a := range actions {
				logger.Info("RenderActions", zap.Any("action", a))
			}

			rendered, err := ExecuteRenderActions[A](
				context.Background(),
				ac.State.Renderer,
				actions,
			)

			logger.Info("Rendered", zap.Any("count", len(rendered)))

			if err != nil {
				logger.Error("Error rendering", zap.Error(err))
				return []RenderedElement{}, err
			}

			for _, r := range rendered {
				logger.Info("Rendered", zap.Any("element", r.renderedKind()))
			}

			return rendered, nil

		},
	}

}

type Handler[S any, A any] struct {
	app        Application[S, A]
	appContext *ApplicationContext[S, A]
	// ChatState  InternalChatState[S, A]
}

func NewHandler[S any, A any](app Application[S, A], tc *TelegramContext) *Handler[S, A] {
	// app.HandleInit(tc)

	appState := app.CreateAppState(tc)

	chatState := InternalChatState[S, A]{
		ChatID:           tc.ChatID,
		AppState:         appState,
		RenderedElements: []RenderedElement{},
		InputHandler:     nil,
		CallbackHandler:  nil,
		Renderer:         app.CreateChatRenderer(tc),
	}

	ac := &ApplicationContext[S, A]{
		App:    &app,
		State:  &chatState,
		Logger: GetLogger().With(zap.Int("chat_id", int(tc.ChatID))),
	}

	res := app.PreRender(ac)

	return &Handler[S, A]{
		app: app,
		appContext: &ApplicationContext[S, A]{
			App:    &app,
			State:  &res.InternalChatState,
			Logger: ac.Logger,
		},
	}
}

func (h *Handler[S, A]) HandleUpdate(tc *TelegramContext) {
	tc.Logger.Info("HandleUpdate")

	if tc.Update.Message != nil && tc.Update.Message.Text != "" {
		h.app.HandleMessage(h.appContext, tc)
	}

	if tc.Update.CallbackQuery != nil {
		h.app.HandleCallback(h.appContext, tc)
	}
}
