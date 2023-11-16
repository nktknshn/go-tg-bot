package main

import (
	tgbot "github.com/nktknshn/go-tg-bot"
)

const (
	ValuePageWelcome  = "Welcome"
	ValuePageMain     = "Main"
	ValuePageSettings = "Settings"
)

type TodoItem struct {
	Text string
	Done bool
	Tags []string
}

type TodoList struct {
	Items []TodoItem
}

type State struct {
	List        TodoList
	CurrentPage string
}

type ActionGoPage struct {
	Page string
}

type Root struct {
	CurrentPage string
	TodoList    TodoList `tgbot:"ctx"`
}

func (a *Root) Render(o tgbot.OO) {

}

type PageSettings struct {
	TodoList TodoList `tgbot:"ctx"`
}

func (a *PageSettings) Render(o tgbot.OO) {

}

type PageMain struct{}

func (a *PageMain) Render(o tgbot.OO) {}

type PageWelcome struct{}

func (a *PageWelcome) Render(o tgbot.OO) {}

var app = tgbot.NewApplication[State, any](
	// initial state
	func(tc *tgbot.TelegramContext) State {
		return State{
			CurrentPage: ValuePageWelcome,
			List:        TodoList{},
		}
	},
	// handle actions
	func(ac *tgbot.ApplicationContext[State, any], tc *tgbot.TelegramContext, a any) {

		switch a := a.(type) {
		case ActionGoPage:
			ac.State.AppState.CurrentPage = a.Page
		}

	},
	// create root component
	func(s State) tgbot.Comp[any] {
		return &Root{
			CurrentPage: s.CurrentPage,
		}
	},
	// create global context
	&tgbot.NewApplicationProps[State, any]{
		CreateGlobalContext: func(ics *tgbot.InternalChatState[State, any]) tgbot.GlobalContext {
			ctx := tgbot.NewGlobalContext()

			ctx.Add("TodoList", ics.AppState.List)

			return ctx
		},
	},
)

func Main() {

}
