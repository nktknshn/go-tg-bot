package tgbot_test

import (
	"reflect"
	"testing"

	tgbot "github.com/nktknshn/go-tg-bot"
)

type App1State struct {
	boolean bool
}

type App1 struct {
	// props
	Counter int
	// local state
	State tgbot.GetSetLocalStateImpl[App1State]
}

func (a App1) Render(o tgbot.OO) {
	ls := a.State.Get(App1State{boolean: false})

	o.Message("Hello")

	if ls.boolean {
		o.Message("World")
	} else {
		o.Message("World2")
	}

	o.Messagef("Counter: %v", a.Counter)
}

func TestRunComponent(t *testing.T) {
	comp := App1{Counter: 1}

	res := tgbot.CreateElements[any](comp, nil)

	t.Logf("res: %v", res)
}

func TestReflect(t *testing.T) {
	app := App1{Counter: 1}

	// getset := tgbot.NewGetSet[any]([]int{0}, nil)
	// tgbot.ReflectCompLocalState[any](app, getset)

	stateField := reflect.ValueOf(app).FieldByName("state")

	t.Logf("stateField: %v", stateField.Type())
}
