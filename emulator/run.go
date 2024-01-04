package emulator

import (
	"github.com/nktknshn/go-tg-bot/tgbot/dispatcher"
	"go.uber.org/zap"
)

func Run(logger *zap.Logger, dispatcher *dispatcher.ChatsDispatcher) {
	bot := NewFakeBot()
	EmulatorMain(bot, dispatcher)
}
