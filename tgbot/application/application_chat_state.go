package application

import (
	"sync"

	"github.com/BooleanCat/go-functional/iter"
	"github.com/gotd/td/tg"
	"github.com/nktknshn/go-tg-bot/tgbot/common"
	"github.com/nktknshn/go-tg-bot/tgbot/component"
	"github.com/nktknshn/go-tg-bot/tgbot/render"
	"go.uber.org/zap"
)

// holds all the state of a chat
type ApplicationChatState[S any, C any] struct {
	ChatID int64

	User *tg.User

	// state of the application
	AppState S

	// local state of the application
	treeState *component.RunResultWithStateTree

	// elements visible to the user
	renderedElements []render.RenderedElement

	// current handler for text messages
	inputHandler common.ChatInputHandler

	// current handler for button presses
	callbackHandler common.ChatCallbackHandler

	// renderer for the messages
	Renderer render.ChatRenderer

	// mutex for locking the state
	lock *sync.Mutex
}

func NewChatState[S any, C any](user *tg.User, appState S) *ApplicationChatState[S, C] {
	return &ApplicationChatState[S, C]{
		ChatID:           user.ID,
		User:             user,
		AppState:         appState,
		renderedElements: []render.RenderedElement{},
		lock:             &sync.Mutex{},
	}
}

func (s *ApplicationChatState[S, C]) ResetRenderedElements() {
	s.renderedElements = make([]render.RenderedElement, 0)
}

func (s *ApplicationChatState[S, C]) LockState(logger *zap.Logger) {
	logger.Debug("LockState")
	s.lock.Lock()
}

func (s *ApplicationChatState[S, C]) UnlockState(logger *zap.Logger) {
	logger.Debug("UnlockState")
	s.lock.Unlock()
}

func (s *ApplicationChatState[S, C]) SetAppState(appState S) {
	s.AppState = appState
}

func (s *ApplicationChatState[S, C]) SetTreeState(treeState *component.RunResultWithStateTree) {
	s.treeState = treeState
}

func (s *ApplicationChatState[S, C]) SetInputHandler(inputHandler common.ChatInputHandler) {
	s.inputHandler = inputHandler
}

func (s *ApplicationChatState[S, C]) SetCallbackHandler(callbackHandler common.ChatCallbackHandler) {
	s.callbackHandler = callbackHandler
}

func (s *ApplicationChatState[S, C]) SetRenderer(renderer render.ChatRenderer) {
	s.Renderer = renderer
}

func (s *ApplicationChatState[S, C]) SetRenderedElements(renderedElements []render.RenderedElement) {
	s.renderedElements = renderedElements
}

func (s *ApplicationChatState[S, C]) RenderedElementsKinds() []string {
	return iter.Map(iter.Lift(s.renderedElements), func(e render.RenderedElement) string {
		return e.RenderedKind()
	}).Collect()
}

func (s *ApplicationChatState[S, C]) RenderedElementsSimpl() []string {
	return iter.Map(iter.Lift(s.renderedElements), func(e render.RenderedElement) string {
		return e.String()
	}).Collect()
}
