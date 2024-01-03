package todo_test

import (
	"path"
	"testing"
	"time"

	"github.com/nktknshn/go-tg-bot/btest"
	emulator "github.com/nktknshn/go-tg-bot/emulator"
	tgbot "github.com/nktknshn/go-tg-bot/tgbot"
	"go.uber.org/zap"

	"github.com/nktknshn/go-tg-bot/emulator/helpers"

	todo "github.com/nktknshn/go-tg-bot/examples/todo/app"
)

func TestTodoApp(t *testing.T) {

	tempDir := t.TempDir()
	usersJson := path.Join(tempDir, "users.json")

	t.Log("TestTodoApp")

	loggers := tgbot.TgbotLoggers{
		Base: tgbot.DevLogger(),
		ChatsDistpatcher: func(l *zap.Logger) *zap.Logger {
			return l.Named("ChatsDistpatcher")
		},
		ChatHandler: func(l *zap.Logger) *zap.Logger {
			return l.Named("ChatHandler")
		},
		Component: func(l *zap.Logger) *zap.Logger {
			// return l.Named("Component")
			return zap.NewNop()
		},
		ApplicationChat: func(l *zap.Logger) *zap.Logger {
			return l.Named("ApplicationChat")
		},
	}

	app := todo.TodoApp(
		todo.TodoAppDeps{
			UserService: todo.NewUserServiceJson(usersJson),
		},
	)

	app.SetLoggers(loggers)

	d := app.ChatsDispatcher()

	d.SetLogger(loggers.ChatsDistpatcher(loggers.Base))

	bot := emulator.NewFakeBot()
	bot.SetDispatcher(d)

	user1 := bot.NewUser(123).SetProfile(emulator.FakeUserInfo{
		Username:  "user1",
		FirstName: "User",
		LastName:  "One",
	})

	user1.SendTextMessage("/start")

	time.Sleep(200 * time.Millisecond)

	println("user1 messages:")

	btest.AssertDisplayedMessages(t, user1, []helpers.MessageSimple{{
		Message: "Welcome User One @user1",
		Buttons: [][]helpers.ButtonSimpl{{{
			Text: "Go to main",
			Data: "Go to main",
		}}}}})

	user1.SendCallbackQuery("Go to main")

	time.Sleep(200 * time.Millisecond)

	user1.SendTextMessage("task 1")
	user1.SendCallbackQuery("Yes")

	time.Sleep(200 * time.Millisecond)

	btest.AssertDisplayedMessages(t, user1, []helpers.MessageSimple{{
		Message: "\n/0 ⭕️ task 1",
		Buttons: nil,
	}})

	d.ResetChats()

	user1.SendTextMessage("/start")

	time.Sleep(200 * time.Millisecond)

	for idx, msg := range user1.DisplayedMessages() {
		println(idx, helpers.MessageAsJson(msg))
	}

	return

	user1.SendTextMessage("task 2")
	user1.SendCallbackQuery("Yes")
	user1.SendTextMessage("/0")
	user1.SendTextMessage("/1")
}