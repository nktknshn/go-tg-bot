package todo_test

import (
	"context"
	"testing"

	emulator "github.com/nktknshn/go-tg-bot/emulator"
	todo "github.com/nktknshn/go-tg-bot/examples/todo/app"
)

func TestTodoApp(t *testing.T) {
	d := todo.TodoApp.ChatsDispatcher()
	bot := emulator.NewFakeBot()

	msg1 := emulator.NewTextMessageUpdate(emulator.TextMessageUpdate{
		Text:        "/start",
		UpdateProps: emulator.UpdateProps{ChatID: 1, UserID: 1},
	})

	bot.AddUserMessage(msg1)
	d.HandleUpdate(context.Background(), bot, msg1)

	d.HandleUpdate(
		context.Background(),
		bot,
		emulator.NewCallbackQueryUpdate(
			emulator.CallbackQueryUpdate{
				Data:        "Go to main",
				UpdateProps: emulator.UpdateProps{ChatID: 1, UserID: 1},
			},
		),
	)

	msg2 := emulator.NewTextMessageUpdate(emulator.TextMessageUpdate{
		Text:        "test",
		UpdateProps: emulator.UpdateProps{ChatID: 1, UserID: 1},
	})

	bot.AddUserMessage(msg2)
	d.HandleUpdate(context.Background(), bot, msg2)
}
