package emulator

import (
	tgbot "github.com/nktknshn/go-tg-bot"
	"go.uber.org/zap"
)

func RunEmulator(logger *zap.Logger, dispatcher *tgbot.ChatsDispatcher) {
	bot := NewFakeBot()
	EmulatorMain(bot, dispatcher)
}
