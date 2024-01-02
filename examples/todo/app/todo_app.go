package todo

import (
	"fmt"

	tgbot "github.com/nktknshn/go-tg-bot"
)

// application state
type TodoState struct {
	Username    string
	List        TodoList
	CurrentPage string
	Error       string
	User        *User
}

// Context shared between all components
type TodoGlobalContext struct {
	Username string
	Settings map[string]string
	TodoList TodoList
}

type App = tgbot.Application[TodoState, TodoGlobalContext]
type AppContext = tgbot.ApplicationContext[TodoState, TodoGlobalContext]

type TodoAppDeps struct {
	UserService UserService
}

func TodoApp(deps TodoAppDeps) *App {
	return tgbot.NewApplication(
		// initial state
		func(tc *tgbot.TelegramContext) TodoState {

			user, err := deps.UserService.GetUser(tc.ChatID)

			if err != nil {
				return TodoState{
					Error: fmt.Sprintf("failed to get user: %v", err),
				}
			}

			if user == nil {
				user = UserFromTgUser(tc.Update.User)

				err := deps.UserService.SaveUser(user)

				if err != nil {
					return TodoState{
						Error: fmt.Sprintf("failed to save user: %v", err),
					}
				}
			}

			username := tgbot.UpdateGetUsername(tc.Update)

			return TodoState{
				User:        user,
				Username:    username,
				CurrentPage: ValuePageWelcome,
				List:        user.TodoList,
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
	).WithCreateGlobalContext(
		func(cs *tgbot.ChatState[TodoState, TodoGlobalContext]) TodoGlobalContext {
			return TodoGlobalContext{
				TodoList: cs.AppState.List,
				Username: cs.AppState.Username,
				Settings: make(map[string]string),
			}
		},
	)
}
