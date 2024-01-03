package telegram

import "github.com/gotd/td/tg"

type BotUpdate struct {
	UpdateClass tg.UpdateClass
	User        *tg.User
	Entities    tg.Entities
}

func (bu BotUpdate) GetNewMessageUpdate() (*tg.UpdateNewMessage, bool) {
	if update, ok := bu.UpdateClass.(*tg.UpdateNewMessage); ok {
		return update, true
	}

	return nil, false
}

func (bu BotUpdate) GetCallbackQueryUpdate() (*tg.UpdateBotCallbackQuery, bool) {
	if update, ok := bu.UpdateClass.(*tg.UpdateBotCallbackQuery); ok {
		return update, true
	}

	return nil, false
}
