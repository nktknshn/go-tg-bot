package todo

import tgbot "github.com/nktknshn/go-tg-bot"

type TodoState struct {
	Username    string
	List        TodoList
	CurrentPage string
}

type TodoGlobalContext struct {
	Username string
	Settings map[string]string
	TodoList TodoList
}

type ApplicationContext = tgbot.ApplicationContext[TodoState, TodoGlobalContext]

var TodoApp = tgbot.NewApplication[TodoState, TodoGlobalContext](
	// initial state
	func(tc *tgbot.TelegramContext) TodoState {

		username := tgbot.UpdateGetUsername(tc.Update)

		return TodoState{
			Username:    username,
			CurrentPage: ValuePageWelcome,
			List:        TodoList{},
		}
	},
	// create root component
	func(s TodoState) tgbot.Comp {
		return &RootComponent{
			CurrentPage: s.CurrentPage,
		}
	},
	// handle actions
	actionsReducer,
	// create global context
	&tgbot.NewApplicationProps[TodoState, TodoGlobalContext]{
		CreateGlobalContext: func(cs *tgbot.ChatState[TodoState, TodoGlobalContext]) TodoGlobalContext {
			return TodoGlobalContext{
				TodoList: cs.AppState.List,
				Username: cs.AppState.Username,
				Settings: make(map[string]string),
			}
		},
	},
)
