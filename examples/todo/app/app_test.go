package todo_test

import (
	"testing"

	emulator "github.com/nktknshn/go-tg-bot/emulator"
	todo "github.com/nktknshn/go-tg-bot/examples/todo/app"
)

func TestTodoApp(t *testing.T) {
	d := todo.TodoApp.ChatsDispatcher()
	bot := emulator.NewFakeBot()
	bot.SetDispatcher(d)

	user1 := bot.NewUser()

	user1.SendTextMessage("/start")

	println("user1 messages:")

	for _, message := range user1.DisplayedMessages() {
		println(message.Message)
	}

	return

	user1.SendCallbackQuery("Go to main")
	user1.SendTextMessage("task 1")
	user1.SendCallbackQuery("Yes")
	user1.SendTextMessage("task 2")
	user1.SendCallbackQuery("Yes")
	user1.SendTextMessage("/0")
	user1.SendTextMessage("/1")
}
