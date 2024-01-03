package tgbot

// Return to move to the next input handler

type ActionNext struct{}

func Next() ActionNext {
	return ActionNext{}
}

func (n ActionNext) String() string {
	return "Next"
}

// reload interface
type ActionReload struct{}

func (a ActionReload) String() string {
	return "ActionReload"
}
