package application

import "github.com/nktknshn/go-tg-bot/tgbot/telegram"

// implements ChatHandler
func (ac *ApplicationChat[S, C]) HandleUpdate(tc *telegram.TelegramUpdateContext) {
	if tcm, ok := tc.AsTextMessage(); ok && tcm.Message.Message != "" {
		ac.App.HandleMessage(ac, tcm)
	} else if tccb, ok := tc.AsCallback(); ok {
		ac.App.HandleCallback(ac, tccb)
	}
}
