package tgbot

import (
	"go.uber.org/zap"
)

// Handles some internal actions sent from handlers
func internalActionHandle[S any, C any](ac *ApplicationChat[S, C], tc *TelegramContext, action any) {
	tc.Logger.Debug("HandleAction", zap.Any("action", reflectStructName(action)))

	switch a := action.(type) {
	case ActionReload:
		ac.State.ResetRenderedElements()
		return
	case actionLocalState[any]:
		tc.Logger.Debug("ActionLocalState was caught. Applying it to the local state tree.",
			zap.Any("index", a.Index),
			zap.Any("LocalStateTree", ac.State.treeState.LocalStateTree),
		)

		ac.State.treeState.LocalStateTree.Set(a.Index[1:], a.F)

		tc.Logger.Debug("Updated LocalStateTree",
			zap.Any("LocalStateTree", ac.State.treeState.LocalStateTree),
		)
		return
	case []any:
		tc.Logger.Debug("A list of actions was caught")

		for _, a := range a {
			internalActionHandle[S, C](ac, tc, a)
		}
		return
	}

	ac.App.HandleAction(ac, tc, action)
}
