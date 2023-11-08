package tgbot

// A is a type of returned Action to be used in actions reducers
type InputHandler[A any] func(*TelegramContext) (A, error)
type CallbackHandler[A any] func(*TelegramContext) (A, error)

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

type RenderFuncType[S any, A any] func(*ApplicationContext[S, A]) []RenderedElement

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
	RenderFunc RenderFuncType[S, A]

	CreateChatRenderer func(*TelegramContext) ChatRenderer
}

func (a *Application[S, A]) NewHandler(tc *TelegramContext) *Handler[S, A] {
	return NewHandler[S, A](*a, tc)
}

func (a *Application[S, A]) PreRender() {

}

type Handler[S any, A any] struct {
	justCreated bool
	ChatState   InternalChatState[S, A]
}

func NewHandler[S any, A any](app Application[S, A], tc *TelegramContext) *Handler[S, A] {
	// app.HandleInit(tc)

	appState := app.CreateAppState(tc)

	app.RenderFunc()

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
