package todo

import tgbot "github.com/nktknshn/go-tg-bot"

type ActionGoPage struct {
	Page string
}

type ActionAddTodoItem struct {
	Text string
}

type ActionMarkDone struct {
	ItemIndex int
}

type ActionItemDelete struct {
	ItemIndex int
}

var actionsReducer = func(ac *ApplicationContext, tc *tgbot.TelegramContext, a any) {
	appState := &ac.State.AppState

	switch a := a.(type) {
	case ActionGoPage:
		appState.CurrentPage = a.Page
	case ActionAddTodoItem:
		appState.List.Items = append(
			appState.List.Items,
			TodoItem{Text: a.Text},
		)
	case ActionMarkDone:
		appState.List.Items[a.ItemIndex].Done = true
	case ActionItemDelete:
		appState.List.Items = append(
			appState.List.Items[:a.ItemIndex],
			appState.List.Items[a.ItemIndex+1:]...,
		)
	}

}
