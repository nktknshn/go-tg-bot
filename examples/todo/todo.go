package main

import (
	"context"
	"flag"

	"github.com/nktknshn/go-tg-bot/emulator"
	todo "github.com/nktknshn/go-tg-bot/examples/todo/app"
	"github.com/nktknshn/go-tg-bot/tgbot/dispatcher"
	"github.com/nktknshn/go-tg-bot/tgbot/gotd"
	"github.com/nktknshn/go-tg-bot/tgbot/logging"
)

func main() {
	flag.Parse()
	runEmul := len(flag.Args()) > 0 && flag.Args()[0] == "e"

	userService := todo.NewUserServiceJson("/tmp/users.json")

	app := todo.TodoApp(todo.TodoAppDeps{
		UserService: userService,
	})

	dispatcher := dispatcher.ForApplication(app)

	if runEmul {
		emulator.Run(
			logging.DevLogger(),
			dispatcher,
		)
	} else {
		gotd.Run(
			context.Background(),
			logging.DevLogger(),
			dispatcher,
		)
	}

}
