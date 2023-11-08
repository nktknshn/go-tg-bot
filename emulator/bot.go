package emulator

import "github.com/go-telegram/bot"

func NewFakeBot() *bot.Bot {
	return &bot.Bot{}
}
