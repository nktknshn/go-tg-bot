package main

import (
	tgbot "github.com/nktknshn/go-tg-bot"
	"github.com/nktknshn/go-tg-bot/emulator"
	todo "github.com/nktknshn/go-tg-bot/examples/todo/app"
)

func main() {

	var runEmul = true
	dispatcher := todo.TodoApp.ChatsDispatcher()

	if runEmul {
		emulator.RunEmulator(
			tgbot.GetLogger(),
			dispatcher,
		)
	} else {
		tgbot.RunReal(
			tgbot.GetLogger(),
			dispatcher,
		)
	}

}
