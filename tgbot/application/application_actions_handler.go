package application

import (
	"github.com/nktknshn/go-tg-bot/tgbot/common"
	"github.com/nktknshn/go-tg-bot/tgbot/component"
	"github.com/nktknshn/go-tg-bot/tgbot/reflection"
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
	"go.uber.org/zap"
)

// Handles some internal actions sent from handlers
func internalActionHandle[S any, C any](ac *ApplicationChat[S, C], tc *telegram.TelegramUpdateContext, action any) {
	logger := ac.Loggers.Action

	logger.Debug("HandleAction", zap.Any("action", reflection.ReflectStructName(action)))

	switch a := action.(type) {
	case common.ActionReload:
		ac.State.ResetRenderedElements()
		return
	case component.ActionLocalState[any]:
		logger.Debug("ActionLocalState was caught. Applying it to the local state tree.",
			zap.Any("index", a.Index),
			zap.Any("LocalStateTree", ac.State.treeState.LocalStateTree),
		)

		ac.State.treeState.LocalStateTree.Set(a.Index[1:], a.F)

		logger.Debug("Updated LocalStateTree",
			zap.Any("LocalStateTree", ac.State.treeState.LocalStateTree),
		)
		return
	case []any:
		logger.Debug("A list of actions was caught")

		for _, a := range a {
			internalActionHandle[S, C](ac, tc, a)
		}
		return
	}

	ac.App.HandleAction(ac, tc, action)
}
