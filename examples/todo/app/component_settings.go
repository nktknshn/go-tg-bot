package todo

import tgbot "github.com/nktknshn/go-tg-bot"

type PageSettings struct {
	Context TodoGlobalContext
}

func (a *PageSettings) Selector() map[string]string {
	return a.Context.Settings
}

func (a *PageSettings) Render(o tgbot.O) {
	// tgbot.ErrMessageNotFound.Error()
	// tgbot.NewGlobalContextTyped[]()
}
