package main

import (
	"regexp"

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

func (app App1) Render(o tgbot.OO) {
	// app.State.S
}

var rexGeneric = regexp.MustCompile("[A-Za-z0-9_]+\\.[A-Za-z0-9_]+\\[(.*)\\]")

// func getAddr(v interface{}) {
// 	reflect.ValueOf(v).Addr()
// }

func main() {
	app := App1{Counter: 1}

	closure := tgbot.LocalStateClosure[any]{
		Initialized: true,
		Value:       App1State{boolean: true},
	}

	localState := tgbot.NewLocalState[any]([]int{0}, &closure)

	tgbot.ReflectCompLocalState[any](app, localState.Getset)

	// fmt.Println("app.State: ", app.State)

	// // s := tgbot.NewGetSet[App1State]([]int{0}, nil)
	// // app.State = s
	// // getset := tgbot.NewGetSet[any]([]int{0}, nil)
	// // tgbot.ReflectCompLocalState[any](app, getset)

	// appStruct := reflect.ValueOf(app)
	// stateField := appStruct.FieldByName("State")

	// t := stateField.Type()

	// fmt.Printf("stateField type: %v\n", stateField.Type())

	// // fmt.Printf("stateField isZero: %v\n", stateField.IsZero())
	// // fmt.Printf("stateField.Elem(): %v\n", stateField.Elem())

	// typeValue := reflect.New(t)

	// fmt.Printf("typeValue type: %v\n", typeValue.Type())

	// // typeValue.Type().

	// // fmt.Printf("typeValue.Elem() type: %v\n", typeValue.Elem().Type())
	// // fmt.Printf("typeValue.Elem() type: %v\n", reflect.ValueOf(typeValue.Elem()))

	// // stateField.Set(typeValue)

}
