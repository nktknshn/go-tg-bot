package tgbot

import (
	"github.com/nktknshn/go-tg-bot/tgbot/common"
	"github.com/nktknshn/go-tg-bot/tgbot/component"
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
)

type Comp = component.Comp
type O = component.O

type TelegramUpdateContext = telegram.TelegramUpdateContext
type ActionReload = common.ActionReload
type ActionNext = common.ActionNext

var Next = common.Next
