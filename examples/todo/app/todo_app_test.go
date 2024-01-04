package todo_test

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/nktknshn/go-tg-bot/btest"
	"github.com/nktknshn/go-tg-bot/emulator"
	"github.com/nktknshn/go-tg-bot/tgbot/dispatcher"
	"github.com/nktknshn/go-tg-bot/tgbot/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/nktknshn/go-tg-bot/emulator/helpers"

	todo "github.com/nktknshn/go-tg-bot/examples/todo/app"
)

type UserLogsProps struct {
	Base       *zap.Logger
	LogsFolder string
}

func UserLogs(props UserLogsProps) *logging.TgbotLoggers {

	// file := path.Join(props.LogsFolder, "tgbot.log")

	// zapcore.AddSync()

	return &logging.TgbotLoggers{
		Base: props.Base,
		ChatsDistpatcher: func(l *zap.Logger) *zap.Logger {
			return l.Named("ChatsDistpatcher")
		},
		ChatHandler: func(l *zap.Logger) *zap.Logger {
			return l.Named("ChatHandler")
		},
		ApplicationChat: func(l *zap.Logger, chatID int64) *zap.Logger {
			// chat specific logger

			return l.Named("ApplicationChat").WithOptions(
				zap.WrapCore(func(core zapcore.Core) zapcore.Core {

					fname := fmt.Sprintf("user_%d.log", chatID)
					chatLogFile := path.Join(props.LogsFolder, fname)

					file, err := os.OpenFile(chatLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

					if err != nil {
						l.Error("failed to open log file", zap.Error(err))
						return core
					}

					return zapcore.NewCore(
						zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
						zapcore.AddSync(zapcore.Lock(file)),
						zap.DebugLevel,
					)
				}),
			)
		},
		Component: func(l *zap.Logger) *zap.Logger {
			return l.Named("Component").WithOptions()
		},
	}
}

type TestScope struct {
	loggers *logging.TgbotLoggers
	tempDir string
	app     *todo.App
	disp    *dispatcher.ChatsDispatcher
	bot     *emulator.FakeBot
	user1   *emulator.FakeBotUser
}

func NewTestScope(t *testing.T) *TestScope {

	tempDir := t.TempDir()

	usersJson := path.Join(tempDir, "users.json")

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stderr),
		zap.DebugLevel,
	)

	baseLogger := zap.New(
		core,
	)

	loggers := UserLogs(UserLogsProps{
		Base:       baseLogger,
		LogsFolder: tempDir,
	})

	app := todo.TodoApp(
		todo.TodoAppDeps{
			UserService: todo.NewUserServiceJson(usersJson),
		},
	)

	app.SetLoggers(loggers)

	d := dispatcher.ForApplication(app)
	d.SetLogger(loggers.ChatsDistpatcher(loggers.Base))

	bot := emulator.NewFakeBot()
	bot.SetDispatcher(d)

	user1 := bot.NewUser(123).SetProfile(
		emulator.FakeUserInfo{
			Username:  "user1",
			FirstName: "User",
			LastName:  "One",
		})

	return &TestScope{
		loggers: loggers,
		tempDir: tempDir,
		app:     app,
		disp:    d,
		bot:     bot,
		user1:   user1,
	}
}

func TestLogs(t *testing.T) {
	scope := NewTestScope(t)
	scope.user1.SendTextMessage("/start")

	time.Sleep(200 * time.Millisecond)

	data, err := os.ReadFile(path.Join(scope.tempDir, "user_123.log"))

	if err != nil {
		t.Fatal(err)
	}

	print("Log content:")
	println(string(data))
}

func TestTodoAppInit(t *testing.T) {
	scope := NewTestScope(t)
	scope.user1.SendTextMessage("/start")

	time.Sleep(200 * time.Millisecond)

}

func TestTodoApp(t *testing.T) {

	scope := NewTestScope(t)

	user1 := scope.user1
	d := scope.disp

	user1.SendTextMessage("/start")

	time.Sleep(200 * time.Millisecond)

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
