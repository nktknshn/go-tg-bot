package tgbot

import (
	"go.uber.org/zap"
)

func internalHandleAction[S any, C any](ac *ApplicationContext[S, C], tc *TelegramContext, a any) {
	tc.Logger.Debug("HandleAction", zap.Any("action", reflectStructName(a)))

	switch a := a.(type) {
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
			internalHandleAction[S, C](ac, tc, a)
		}
		return
	}

	ac.App.HandleAction(ac, tc, a)
}
