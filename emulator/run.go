package emulator

import (
	tgbot "github.com/nktknshn/go-tg-bot/tgbot"
	"go.uber.org/zap"
)

func Run(logger *zap.Logger, dispatcher *tgbot.ChatsDispatcher) {
	bot := NewFakeBot()
	EmulatorMain(bot, dispatcher)
}
