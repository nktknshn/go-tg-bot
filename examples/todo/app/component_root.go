package todo

import tgbot "github.com/nktknshn/go-tg-bot/tgbot"

const (
	ValuePageWelcome  = "Welcome"
	ValuePageMain     = "Main"
	ValuePageSettings = "Settings"
)

type PageWelcome struct {
	Context TodoGlobalContext
}

func (a *PageWelcome) Selector() string {
	return a.Context.Username
}

func (a *PageWelcome) Render(o tgbot.O) {
	o.Messagef("Welcome %v", a.Selector())
	o.Button("Go to main", func() any {
		return ActionGoPage{Page: ValuePageMain}
	})
}

type RootComponent struct {
	CurrentPage string
}

func (a *RootComponent) Render(o tgbot.O) {
	o.InputHandler(func(s string) any {
		if s == "/start" {
			return &tgbot.ActionReload{}
		}

		return tgbot.ActionNext{}
	})

	switch a.CurrentPage {
	case ValuePageWelcome:
		o.Comp(&PageWelcome{})
	case ValuePageMain:
		o.Comp(&PageTodoList{})
	case ValuePageSettings:
		o.Comp(&PageSettings{})
	}
}
