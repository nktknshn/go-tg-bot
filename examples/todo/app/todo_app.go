package todo

import (
	"fmt"

	tgbot "github.com/nktknshn/go-tg-bot/tgbot"
	"github.com/nktknshn/go-tg-bot/tgbot/application"
	"go.uber.org/zap"
)

// application state
type TodoState struct {
	Username    string
	List        TodoList
	CurrentPage string
	Error       string
	User        *TodoUser
}

// Context shared between all components
type TodoGlobalContext struct {
	Username string
	Settings map[string]string
	TodoList TodoList
}

type TodoAppDeps struct {
	UserService UserService
}

type App = application.Application[TodoState, TodoGlobalContext]
type AppChat = application.ApplicationChat[TodoState, TodoGlobalContext]

func createAppState(app *App, deps TodoAppDeps, tc *tgbot.TelegramUpdateContext, logger *zap.Logger) TodoState {

	logger.Debug("Fetching user")

	user, err := deps.UserService.GetUser(tc.ChatID)

	if err != nil {
		logger.Error("Fetching user", zap.Error(err))

		return TodoState{
			Error: fmt.Sprintf("failed to get user: %v", err),
		}
	}

	if user == nil {
		logger.Debug("User not found. Creating a new one")

		user = TodoUserFromTgUser(tc.Update.User)

		err := deps.UserService.SaveUser(user)

		if err != nil {
			return TodoState{
				Error: fmt.Sprintf("failed to save user: %v", err),
			}
		}
	}

	username := fmt.Sprintf("%v %v @%v", user.FirstName, user.LastName, user.Username)

	return TodoState{
		User:        user,
		Username:    username,
		CurrentPage: ValuePageWelcome,
		List:        user.TodoList,
	}
}

func TodoApp(deps TodoAppDeps) *App {
	return application.New(
		// initial state
		func(app *App, tc *tgbot.TelegramUpdateContext, logger *zap.Logger) TodoState {
			return createAppState(app, deps, tc, logger)
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
	).WithGlobalContext(
		func(cs *application.ApplicationChatState[TodoState, TodoGlobalContext]) TodoGlobalContext {
			return TodoGlobalContext{
				TodoList: cs.AppState.List,
				Username: cs.AppState.Username,
				Settings: make(map[string]string),
			}
		},
	)
}
