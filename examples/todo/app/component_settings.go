package todo

import tgbot "github.com/nktknshn/go-tg-bot"

type PageSettings struct {
	Context AppGlobalContext
}

func (a *PageSettings) Selector() map[string]string {
	return a.Context.Settings
}

func (a *PageSettings) Render(o tgbot.OO) {

}
