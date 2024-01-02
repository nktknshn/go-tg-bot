package todo

import (
	tgbot "github.com/nktknshn/go-tg-bot"
	"go.uber.org/zap"
)

type ActionGoPage struct {
	Page string
}

type ActionAddTodoItem struct {
	Text string
}

type ActionToggle struct {
	ItemIndex int
}

type ActionItemDelete struct {
	ItemIndex int
}

var actionsReducer = func(ac *AppContext, tc *tgbot.TelegramContext, a any) {
	appState := &ac.State.AppState

	switch a := a.(type) {
	case ActionGoPage:
		appState.CurrentPage = a.Page

	case ActionAddTodoItem:
		ac.Logger.Info("adding item",
			zap.String("text", a.Text),
			zap.Int("len", len(appState.List.Items)),
		)

		appState.List.Items = append(
			appState.List.Items,
			TodoItem{Text: a.Text},
		)
	case ActionToggle:
		appState.List.Items[a.ItemIndex].Done = !appState.List.Items[a.ItemIndex].Done
	case ActionItemDelete:
		appState.List.Items = append(
			appState.List.Items[:a.ItemIndex],
			appState.List.Items[a.ItemIndex+1:]...,
		)
	}

}
