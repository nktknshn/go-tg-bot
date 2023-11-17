package main

import (
	"flag"

	tgbot "github.com/nktknshn/go-tg-bot"
	"github.com/nktknshn/go-tg-bot/emulator"
	todo "github.com/nktknshn/go-tg-bot/examples/todo/app"
)

func main() {
	flag.Parse()
	runEmul := len(flag.Args()) > 0 && flag.Args()[0] == "e"

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
