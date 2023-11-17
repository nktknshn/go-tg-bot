package tgbot

import "go.uber.org/zap"

// Next handler
type Next struct{}

func (n Next) String() string {
	return "Next"
}

// reload interface
type ActionReload struct{}

func (a ActionReload) String() string {
	return "ActionReload"
}

func internalHandleAction[S any, A any, C any](ac *ApplicationContext[S, A, C], tc *TelegramContext, a any) {
	tc.Logger.Debug("HandleAction", zap.Any("action", ReflectStructName(a)))

	switch a := a.(type) {
	case ActionReload:
		ac.State.RenderedElements = make([]RenderedElement, 0)
	case ActionLocalState[any]:
		tc.Logger.Debug("ActionLocalState was caught. Applying it to the local state tree.", zap.Any("index", a.index))
		ac.State.TreeState.LocalStateTree.Set(a.index, a.f)
		return
	case []A:
		tc.Logger.Debug("A list of actions was caught")

		for _, a := range a {
			internalHandleAction(ac, tc, a)
		}
	}

	ac.App.HandleAction(ac, tc, a)
}
