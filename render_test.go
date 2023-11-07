package tgbot_test

type AppProps struct {
	Name    string
	Counter int
}

// func HelloMessage(name string) tgbot.Comp {
// 	return func(o tgbot.O) {
// 		o.Message(fmt.Sprintf("Hello, %v", name))
// 	}
// }

// func Counter(value int) tgbot.Comp {
// 	return func(o tgbot.O) {
// 		messageText := fmt.Sprintf("Counter %v", value)

// 		o.Message(messageText)
// 		o.Button("add", func() {})
// 		o.Button("sub", func() {})

// 	}
// }

// func App(props *AppProps) tgbot.Comp {
// 	return func(o tgbot.O) {
// 		o.Comp(HelloMessage(props.Name))
// 		o.Comp(Counter(props.Counter))
// 	}
// }

// func TestRenderComp(t *testing.T) {
// 	tgbot.ComponentToElements(
// 		App,
// 		&AppProps{Counter: 1},
// 	)

// }
