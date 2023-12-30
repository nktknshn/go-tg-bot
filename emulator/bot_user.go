package emulator

import (
	"context"

	"github.com/gotd/td/tg"
	tgbot "github.com/nktknshn/go-tg-bot"
)

type FakeBotUser struct {
	UserID int64
	ChatID int64
	Bot    *FakeBot
}

func (fbu *FakeBotUser) GetUser() *tg.User {
	return &tg.User{
		ID:       fbu.UserID,
		Username: "username",
	}
}

// construct BotUpdate
func (u *FakeBotUser) SendTextMessage(text string) tg.UpdateClass {

	bu := tgbot.BotUpdate{}

	textMessage := &tg.Message{
		Message: text,
	}

	bu.UpdateClass = &tg.UpdateNewMessage{
		Message: textMessage,
	}

	bu.User = u.GetUser()

	u.Bot.dispatcher.HandleUpdate(context.Background(), u.Bot, bu)

	return bu.UpdateClass
}

func (u *FakeBotUser) SendCallbackQuery(data string) tg.UpdateClass {
	bu := tgbot.BotUpdate{}

	bu.UpdateClass = &tg.UpdateBotCallbackQuery{
		Data: []byte(data),
	}

	bu.User = u.GetUser()

	u.Bot.dispatcher.HandleUpdate(context.Background(), u.Bot, bu)

	return bu.UpdateClass
}