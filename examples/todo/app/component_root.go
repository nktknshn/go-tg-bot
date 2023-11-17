package todo

import tgbot "github.com/nktknshn/go-tg-bot"

const (
	ValuePageWelcome  = "Welcome"
	ValuePageMain     = "Main"
	ValuePageSettings = "Settings"
)

type PageWelcome struct {
	Context AppGlobalContext
}

func (a *PageWelcome) Selector() string {
	return a.Context.Username
}

func (a *PageWelcome) Render(o tgbot.OO) {
	o.Messagef("Welcome %v", a.Selector())
	o.Button("Go to main", func() any {
		return ActionGoPage{Page: ValuePageMain}
	})
}

type RootComponent struct {
	CurrentPage string
}

func (a *RootComponent) Render(o tgbot.OO) {
	o.InputHandler(func(s string) any {
		if s == "/start" {
			return &tgbot.ActionReload{}
		}

		return tgbot.Next{}
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
