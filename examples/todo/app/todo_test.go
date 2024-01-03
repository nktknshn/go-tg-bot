package todo

import (
	"testing"

	"github.com/nktknshn/go-tg-bot/emulator"
	tgbot "github.com/nktknshn/go-tg-bot/tgbot"
)

func TestTodo(t *testing.T) {
	_ = tgbot.ActionReload{}

	userService := NewUserServiceJson("/tmp/users.json")

	app := TodoApp(TodoAppDeps{
		userService,
	})

	dispatcher := app.ChatsDispatcher()
	bot := emulator.NewFakeBot()

	bot.SetDispatcher(dispatcher)

	user1 := bot.NewUser(emulator.RandomUserID())

	user1.SendTextMessage("/start")

	// bot.Test3()
}
