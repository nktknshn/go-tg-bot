package tgbot

import (
	"go.uber.org/zap"
)

// Handles some internal actions sent from handlers
func internalActionHandle[S any, C any](ac *ApplicationChat[S, C], tc *TelegramUpdateContext, action any, logger *zap.Logger) {
	logger.Debug("HandleAction", zap.Any("action", reflectStructName(action)))

	switch a := action.(type) {
	case ActionReload:
		ac.State.ResetRenderedElements()
		return
	case actionLocalState[any]:
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
			internalActionHandle[S, C](ac, tc, a, logger)
		}
		return
	}

	ac.App.HandleAction(ac, tc, action)
}
