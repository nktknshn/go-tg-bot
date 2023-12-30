package todo

import (
	"testing"

	tgbot "github.com/nktknshn/go-tg-bot"
	"github.com/nktknshn/go-tg-bot/emulator"
)

func TestTodo(t *testing.T) {
	_ = tgbot.ActionReload{}

	dispatcher := TodoApp.ChatsDispatcher()
	bot := emulator.NewFakeBot()

	bot.SetDispatcher(dispatcher)

	user1 := bot.NewUser()

	user1.SendTextMessage("/start")

	// bot.Test3()
}
