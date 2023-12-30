package tgbot

import (
	"sync"

	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

type ChatState[S any, C any] struct {
	ChatID int64

	User *tg.User

	// state of the application
	AppState S

	// state of the application
	treeState *runResultWithStateTree

	// elements visible to the user
	renderedElements []RenderedElement

	// handler for text messages
	inputHandler chatInputHandler

	// handler for callback queries
	callbackHandler chatCallbackHandler

	Renderer ChatRenderer

	lock *sync.Mutex
}

func NewChatState[S any, C any](user *tg.User, appState S) *ChatState[S, C] {
	return &ChatState[S, C]{
		ChatID:           user.ID,
		User:             user,
		AppState:         appState,
		renderedElements: []RenderedElement{},
		lock:             &sync.Mutex{},
	}
}

func (s *ChatState[S, C]) ResetRenderedElements() {
	s.renderedElements = make([]RenderedElement, 0)
}

func (s *ChatState[S, C]) LockState(logger *zap.Logger) {
	logger.Debug("LockState")
	s.lock.Lock()
}

func (s *ChatState[S, C]) UnlockState(logger *zap.Logger) {
	logger.Debug("UnlockState")
	s.lock.Unlock()
}
