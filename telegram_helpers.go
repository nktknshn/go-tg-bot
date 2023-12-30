package tgbot

import (
	"fmt"
	"io"
	"os"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tg"
)

func TryReadBotTokenFile() (string, error) {
	res, err := os.ReadFile("../bot.txt")

	if err != nil {
		return "", err
	}

	return string(res), nil
}

func UpdateGetUsername(update BotUpdate) string {

	return fmt.Sprintf("%v %v %v", update.User.FirstName, update.User.LastName, update.User.Username)
}

func UpdateToString(update tg.UpdateClass) string {
	// result := ""

	// if update.Message != nil {
	// 	result += "Message: " + update.Message.Text
	// }

	// if update.CallbackQuery != nil {
	// 	result += "CallbackQuery: " + update.CallbackQuery.Data
	// }

	// if update.InlineQuery != nil {
	// 	result += "InlineQuery: " + update.InlineQuery.Query
	// }

	// if update.ChosenInlineResult != nil {
	// 	result += "ChosenInlineResult: " + update.ChosenInlineResult.Query
	// }

	// if update.ShippingQuery != nil {
	// 	result += "ShippingQuery: " + update.ShippingQuery.InvoicePayload
	// }

	// if update.EditedMessage != nil {
	// 	result += "EditedMessage: " + update.EditedMessage.Text
	// }

	return update.String()
}

func RandInt64(randSource io.Reader) (int64, error) {
	var buf [bin.Word * 2]byte
	if _, err := io.ReadFull(randSource, buf[:]); err != nil {
		return 0, err
	}
	b := &bin.Buffer{Buf: buf[:]}
	return b.Long()
}
