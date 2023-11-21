package tgbot

type AppProps struct {
	Name    string
	Counter int
}

// func HelloMessage(name string) Comp {
// 	return func(o O) {
// 		o.Message(fmt.Sprintf("Hello, %v", name))
// 	}
// }

// func Counter(value int) Comp {
// 	return func(o O) {
// 		messageText := fmt.Sprintf("Counter %v", value)

// 		o.Message(messageText)
// 		o.Button("add", func() {})
// 		o.Button("sub", func() {})

// 	}
// }

// func App(props *AppProps) Comp {
// 	return func(o O) {
// 		o.Comp(HelloMessage(props.Name))
// 		o.Comp(Counter(props.Counter))
// 	}
// }

// func TestRenderComp(t *testing.T) {
// 	ComponentToElements(
// 		App,
// 		&AppProps{Counter: 1},
// 	)

// }
