package tgbot

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// A is a type of returned Action to be used in actions reducers
type InputHandler[A any] func(*TelegramContext) (A, error)
type CallbackHandler[A any] func(*TelegramContext) (*A, error)

type InternalChatState[S any, A any] struct {
	ChatID int64
	// state of the application
	AppState S

	// elements visible to the user
	RenderedElements []RenderedElement

	// handler for text messages
	InputHandler InputHandler[A]

	// handler for callback queries
	CallbackHandler CallbackHandler[A]

	Renderer ChatRenderer
}

// func (s *InternalChatState[S]) ModifyState() {
// 	return
// }

type renderFunc[S any] func(S) *renderFuncResult[S]

type renderFuncResult[S any] struct {
	chatState *S
}

type HandleMessageFunc[S any, A any] func(*ApplicationContext[S, A], *TelegramContext)
type HandleCallbackFunc[S any, A any] func(*ApplicationContext[S, A], *TelegramContext)
type HandleInitFunc[S any] func(*TelegramContext)

type HandleActionFunc[S any, A any] func(*ApplicationContext[S, A], *TelegramContext, A)

// type ReducerFuncType[A any, S any] func(InternalChatState[S]) InternalChatState[S]

type RenderFuncType[S any, A any] func(*ApplicationContext[S, A]) error
type StateToCompFuncType[S any, A any] func(S) Comp[A]

type ApplicationContext[S any, A any] struct {
	App   *Application[S, A]
	State *InternalChatState[S, A]
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
					logger.Info("Setting callback handler", zap.String("key", butt.Action))

					callbackHandlers[butt.Action] = func() *A {
						v := butt.OnClick()
						return &v
					}
				}
			}
		}
	}

	return callbackHandlers
}

func callbackMapToHandler[A any](cbmap map[string]func() *A) CallbackHandler[A] {
	return func(tc *TelegramContext) (*A, error) {

		logger.Info("Callback handler", zap.String("data", tc.Update.CallbackQuery.Data))

		key := tc.Update.CallbackQuery.Data

		if handler, ok := cbmap[key]; ok {
			// ac.App.HandleAction(ac, tc, )
			logger.Info("Calling handler", zap.String("data", tc.Update.CallbackQuery.Data))

			return handler(), nil
		} else {
			err := fmt.Errorf("no handler for callback %v", key)
			logger.Error("No handler for callback", zap.String("key", key))
			return nil, err
		}

	}
}

func (a *Application[S, A]) PreRender(ac *ApplicationContext[S, A]) (*PreRenderData[S, A], error) {

	els := ComponentToElements2(
		ac.App.StateToComp(ac.State.AppState),
	)

	res, err := ElementsToMessagesAndHandlers[A](els)

	if err != nil {
		return nil, err
	}

	var inputHandler InputHandler[A]
	callbackMap := getCallbackHandlersMap[A](res.OutcomingMessages)
	callbackHandler := callbackMapToHandler[A](callbackMap)

	if len(res.InputHandlers) > 0 {
		inputHandler = res.InputHandlers[0].Handler
	}

	stateCopy := InternalChatState[S, A]{
		ChatID:           ac.State.ChatID,
		AppState:         ac.State.AppState,
		RenderedElements: ac.State.RenderedElements,
		InputHandler:     inputHandler,
		CallbackHandler:  callbackHandler,
		Renderer:         ac.State.Renderer,
	}

	return &PreRenderData[S, A]{
		InternalChatState: stateCopy,
		ExecuteRender: func(renderer ChatRenderer) ([]RenderedElement, error) {
			actions := GetRenderActions(
				ac.State.RenderedElements,
				res.OutcomingMessages,
			)

			logger.Info("RenderActions", zap.Any("count", len(actions)))

			rendered, err := ExecuteRenderActions(
				context.Background(),
				ac.State.Renderer,
				actions,
			)

			logger.Info("Rendered", zap.Any("count", len(rendered)))

			if err != nil {
				logger.Error("Error rendering", zap.Error(err))
				return []RenderedElement{}, err
			}

			return rendered, nil

		},
	}, nil

}

type Handler[S any, A any] struct {
	justCreated bool
	ChatState   InternalChatState[S, A]
}

func NewHandler[S any, A any](app Application[S, A], tc *TelegramContext) *Handler[S, A] {
	// app.HandleInit(tc)

	appState := app.CreateAppState(tc)

	// app.RenderFunc()

	return &Handler[S, A]{
		justCreated: true,
		ChatState: InternalChatState[S, A]{
			ChatID:           tc.ChatID,
			AppState:         appState,
			RenderedElements: []RenderedElement{},
			// InputHandler: func(chc *ChatHandlerContext, u *models.Update) A {
			// 	return 0
			// },
			// CallbackHandler: func(chc *ChatHandlerContext, u *models.Update) A {
			// 	return 0
			// },
			Renderer: app.CreateChatRenderer(tc),
		},
	}
}

func (h *Handler[S, A]) HandleUpdate(tc *TelegramContext) {
	tc.Logger.Info("HandleUpdate")

	if tc.Update.Message != nil {
		if h.ChatState.InputHandler == nil {
			tc.Logger.Info("No input handler, skipping")
			return
		}

		h.ChatState.InputHandler(tc)
		return
	}

	if tc.Update.CallbackQuery != nil {
		if h.ChatState.CallbackHandler == nil {
			tc.Logger.Info("No callback handler, skipping")
			return
		}
		h.ChatState.CallbackHandler(tc)
		return
	}
}
