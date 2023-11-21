package tgbot

import (
	"os"

	"github.com/go-telegram/bot/models"
)

func TryReadBotTokenFile() (string, error) {
	res, err := os.ReadFile("../bot.txt")

	if err != nil {
		return "", err
	}

	return string(res), nil
}

func UpdateGetUsername(update *models.Update) string {
	username := ""

	if update.Message != nil {
		username = update.Message.From.Username
	}

	return username
}

func getUpdateChatId(update *models.Update) int64 {

	if update.Message != nil {
		return update.Message.Chat.ID
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	}

	if update.InlineQuery != nil {
		return update.InlineQuery.From.ID
	}

	if update.ChosenInlineResult != nil {
		return update.ChosenInlineResult.From.ID
	}

	if update.ShippingQuery != nil {
		return update.ShippingQuery.From.ID
	}

	if update.EditedMessage != nil {
		return update.EditedMessage.Chat.ID
	}

	return 0
}

func UpdateToString(update *models.Update) string {
	result := ""

	if update.Message != nil {
		result += "Message: " + update.Message.Text
	}

	if update.CallbackQuery != nil {
		result += "CallbackQuery: " + update.CallbackQuery.Data
	}

	if update.InlineQuery != nil {
		result += "InlineQuery: " + update.InlineQuery.Query
	}

	if update.ChosenInlineResult != nil {
		result += "ChosenInlineResult: " + update.ChosenInlineResult.Query
	}

	if update.ShippingQuery != nil {
		result += "ShippingQuery: " + update.ShippingQuery.InvoicePayload
	}

	if update.EditedMessage != nil {
		result += "EditedMessage: " + update.EditedMessage.Text
	}

	return result
}
