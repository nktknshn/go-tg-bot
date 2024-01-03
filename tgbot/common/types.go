package common

type CallbackResult struct {
	Action any

	// do not answer to the callback
	NoAnswer bool
}

type ChatCallbackHandler func(string) *CallbackResult

// Returns Next if no action is needed
type ChatInputHandler func(string) any
