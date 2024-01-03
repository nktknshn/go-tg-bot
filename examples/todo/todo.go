package main

import (
	"context"
	"flag"

	"github.com/nktknshn/go-tg-bot/emulator"
	todo "github.com/nktknshn/go-tg-bot/examples/todo/app"
	tgbot "github.com/nktknshn/go-tg-bot/tgbot"
)

func main() {
	flag.Parse()
	runEmul := len(flag.Args()) > 0 && flag.Args()[0] == "e"

	userService := todo.NewUserServiceJson("/tmp/users.json")

	app := todo.TodoApp(todo.TodoAppDeps{
		UserService: userService,
	})

	dispatcher := app.ChatsDispatcher()

	if runEmul {
		emulator.Run(
			tgbot.DevLogger(),
			dispatcher,
		)
	} else {
		tgbot.Run(
			context.Background(),
			tgbot.DevLogger(),
			dispatcher,
		)
	}

}
