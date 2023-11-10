package main

import (
	tgbot "github.com/nktknshn/go-tg-bot"
)

type C[P any] func(P, tgbot.Z)

func WelcomZ(usename string, z Z) {
	var locs tgbot.GetSetLocalState[WelcomState, Action]

	if locs.Get().hideName {
		z.Message("Welcome")
	} else {
		z.MessagePartf("Welcome %v", usename)
		z.MessagePart("/hide_name to hide your name")
		z.MessageComplete()
	}

}

type Z interface {
	tgbot.O[any]
	CompZ(func(z Z))
	Comp3(Comp3)
	C(Comp3)
}

func A[P any](cons func(P, Z), props P) func(z Z) {
	return func(z Z) {
		cons(props, z)
	}
}

func App2(props Props, z Z) {
	WelcomZ(props.username, z)

	// z.CompZ(A(WelcomZ, props.username))

	z.Messagef("Counter value: %v", props.counter)

	if props.err != nil {
		z.Messagef("Error: %v", props.err)
	}

	z.Button("Increment", func() any {
		return ActionCounter{Increment: 1}
	})
	z.Button("Decrement", func() any {
		return ActionCounter{Increment: -1}
	})
}

type WelcomZ2 struct {
	username string
	state    tgbot.LocalStateIniter[WelcomState]
}

func (w *WelcomZ2) Render(z Z) {
	ls := w.state.Init(WelcomState{})

	if ls.Get().hideName {

	}

	z.InputHandler(func(s string) any {
		if s == "/hide_name" {
			return ls.Set(func(s WelcomState) {
				s.hideName = true
			})
		}

		return nil
	})

	z.MessagePartf("Welcome username")
	z.MessagePart("/hide_name to hide your name")
	z.MessageComplete()
}

type App3 struct {
	props Props
}

type Comp3 interface {
	Render(z Z)
}

func (a *App3) Render(z Z) {
	// locs := z.LocalState(WelcomState{})
	// z.CompZ(A(WelcomZ, props.username))

	z.C(&WelcomZ2{username: a.props.username})

	z.Messagef("Counter value: %v", a.props.counter)

	if a.props.err != nil {
		z.Messagef("Error: %v", a.props.err)
	}

	z.Button("Increment", func() any {
		return ActionCounter{Increment: 1}
	})
	z.Button("Decrement", func() any {
		return ActionCounter{Increment: -1}
	})
}
