package emulator

/*
Let's emulate
*/

type Butt struct {
	Title          string
	CallbackString string
}

type ButtonsRow struct {
	Butts []Butt
}

type BottomButtons struct {
	Butts []string
}

type ButtonsBotton struct {
}

type MessageButtons struct {
	Rows []ButtonsRow
}

type MessageBox struct {
	Text    string
	Buttons MessageButtons
}

type Menu struct {
	Items []string
}

type DrawInput struct {
	Boxes         []MessageBox
	BottomButtons []BottomButtons
	Menu          Menu
}

type CallbackHandlers func(string)
type UserInputHandler func(string)

type ActionsHandler struct {
	// handlers for callback buttons
	CallbackHandlers
	// handlers for commands, user input, text buttons
	UserInputHandler
}
