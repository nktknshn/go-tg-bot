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

type ChatState[S any, C any] struct {
	ChatID int64

	User *tg.User

	// state of the application
	AppState S

	// state of the application
	treeState *component.RunResultWithStateTree

	// elements visible to the user
	renderedElements []render.RenderedElement

	// handler for text messages
	inputHandler common.ChatInputHandler

	// handler for callback queries
	callbackHandler common.ChatCallbackHandler

	Renderer render.ChatRenderer

	lock *sync.Mutex
}

func NewChatState[S any, C any](user *tg.User, appState S) *ChatState[S, C] {
	return &ChatState[S, C]{
		ChatID:           user.ID,
		User:             user,
		AppState:         appState,
		renderedElements: []render.RenderedElement{},
		lock:             &sync.Mutex{},
	}
}

func (s *ChatState[S, C]) ResetRenderedElements() {
	s.renderedElements = make([]render.RenderedElement, 0)
}

func (s *ChatState[S, C]) LockState(logger *zap.Logger) {
	logger.Debug("LockState")
	s.lock.Lock()
}

func (s *ChatState[S, C]) UnlockState(logger *zap.Logger) {
	logger.Debug("UnlockState")
	s.lock.Unlock()
}

func (s *ChatState[S, C]) SetAppState(appState S) {
	s.AppState = appState
}

func (s *ChatState[S, C]) SetTreeState(treeState *component.RunResultWithStateTree) {
	s.treeState = treeState
}

func (s *ChatState[S, C]) SetInputHandler(inputHandler common.ChatInputHandler) {
	s.inputHandler = inputHandler
}

func (s *ChatState[S, C]) SetCallbackHandler(callbackHandler common.ChatCallbackHandler) {
	s.callbackHandler = callbackHandler
}

func (s *ChatState[S, C]) SetRenderer(renderer render.ChatRenderer) {
	s.Renderer = renderer
}

func (s *ChatState[S, C]) SetRenderedElements(renderedElements []render.RenderedElement) {
	s.renderedElements = renderedElements
}

func (s *ChatState[S, C]) RenderedElementsKinds() []string {
	return iter.Map(iter.Lift(s.renderedElements), func(e render.RenderedElement) string {
		return e.RenderedKind()
	}).Collect()
}

func (s *ChatState[S, C]) RenderedElementsSimpl() []string {
	return iter.Map(iter.Lift(s.renderedElements), func(e render.RenderedElement) string {
		return e.String()
	}).Collect()
}