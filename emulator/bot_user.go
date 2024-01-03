package emulator

import (
	"context"
	"math/rand"

	"github.com/gotd/td/tg"
	tgbot "github.com/nktknshn/go-tg-bot/tgbot"
)

type FakeUserInfo struct {
	Username   string
	FirstName  string
	LastName   string
	AccessHash int64
}

type FakeBotUser struct {
	UserID int64
	Bot    *FakeBot

	AccessHash int64
	Username   string
	FistName   string
	LastName   string
}

func RandomUserID() int64 {
	return rand.Int63()
}

func NewFakeBotUser(userID int64) *FakeBotUser {
	return &FakeBotUser{
		UserID: userID,
	}
}

func (fbu *FakeBotUser) SetProfile(profie FakeUserInfo) *FakeBotUser {
	fbu.Username = profie.Username
	fbu.FistName = profie.FirstName
	fbu.LastName = profie.LastName
	fbu.AccessHash = profie.AccessHash

	return fbu
}

func (fbu *FakeBotUser) GetTgUser() *tg.User {
	return &tg.User{
		ID:         fbu.UserID,
		Username:   fbu.Username,
		FirstName:  fbu.FistName,
		LastName:   fbu.LastName,
		AccessHash: fbu.AccessHash,
	}
}

func (u *FakeBotUser) DisplayedMessages() []*tg.Message {
	return u.Bot.DisplayedMessages(u.UserID)
}

func (u *FakeBotUser) SendTextMessage(text string) tg.UpdateClass {

	bu := tgbot.BotUpdate{}

	textMessage := &tg.Message{
		Message: text,
		PeerID:  &tg.PeerUser{UserID: u.UserID},
	}

	bu.UpdateClass = &tg.UpdateNewMessage{
		Message: textMessage,
	}

	u.Bot.AddUserMessage(textMessage)

	bu.User = u.GetTgUser()

	u.Bot.dispatcher.HandleUpdate(context.Background(), u.Bot, bu)

	return bu.UpdateClass
}

func (u *FakeBotUser) SendCallbackQuery(data string) tg.UpdateClass {
	bu := tgbot.BotUpdate{}

	bu.UpdateClass = &tg.UpdateBotCallbackQuery{
		Data: []byte(data),
	}

	bu.User = u.GetTgUser()

	u.Bot.dispatcher.HandleUpdate(context.Background(), u.Bot, bu)

	return bu.UpdateClass
}
