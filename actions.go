package tgbot

// Return to move to the next input handler
type Next struct{}

func (n Next) String() string {
	return "Next"
}

// reload interface
type ActionReload struct{}

func (a ActionReload) String() string {
	return "ActionReload"
}
