package todo

import tgbot "github.com/nktknshn/go-tg-bot"

type TodoState struct {
	Username    string
	List        TodoList
	CurrentPage string
}

type AppGlobalContext struct {
	Username string
	Settings map[string]string
	TodoList TodoList
}

type ApplicationContext = tgbot.ApplicationContext[TodoState, any, AppGlobalContext]

var TodoApp = tgbot.NewApplication[TodoState, any](
	// initial state
	func(tc *tgbot.TelegramContext) TodoState {
		return TodoState{
			Username:    tc.Update.Message.From.Username,
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
	&tgbot.NewApplicationProps[TodoState, any, AppGlobalContext]{
		CreateGlobalContext: func(ics *tgbot.InternalChatState[TodoState, any, AppGlobalContext]) tgbot.GlobalContextTyped[AppGlobalContext] {

			return tgbot.NewGlobalContextTyped(AppGlobalContext{
				TodoList: ics.AppState.List,
				Username: ics.AppState.Username,
				Settings: make(map[string]string),
			})

		},
	},
)
