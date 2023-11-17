package todo

import tgbot "github.com/nktknshn/go-tg-bot"

type ActionGoPage struct {
	Page string
}

type ActionAddTodoItem struct {
	Text string
}

var actionsReducer = func(ac *ApplicationContext, tc *tgbot.TelegramContext, a any) {
	appState := ac.State.AppState

	switch a := a.(type) {
	case ActionGoPage:
		appState.CurrentPage = a.Page
	case ActionAddTodoItem:
		appState.List.Items = append(
			appState.List.Items,
			TodoItem{Text: a.Text},
		)
	}

}
