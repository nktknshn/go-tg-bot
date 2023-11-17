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

type ApplicationContext = tgbot.ApplicationContext[TodoState, any, TodoGlobalContext]

var TodoApp = tgbot.NewApplication[TodoState, any](
	// initial state
	func(tc *tgbot.TelegramContext) TodoState {

		username := tgbot.GetUsername(tc.Update)

		return TodoState{
			Username:    username,
			CurrentPage: ValuePageWelcome,
			List:        TodoList{},
		}
	},
	// create root component
	func(s TodoState) tgbot.Comp[any] {
		return &RootComponent{
			CurrentPage: s.CurrentPage,
		}
	},
	// handle actions
	actionsReducer,
	// create global context
	&tgbot.NewApplicationProps[TodoState, any, TodoGlobalContext]{
		CreateGlobalContext: func(ics *tgbot.InternalChatState[TodoState, any, TodoGlobalContext]) tgbot.GlobalContextTyped[TodoGlobalContext] {

			return tgbot.NewGlobalContextTyped(TodoGlobalContext{
				TodoList: ics.AppState.List,
				Username: ics.AppState.Username,
				Settings: make(map[string]string),
			})

		},
	},
)
