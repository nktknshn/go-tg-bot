package helpers

import (
	"os"
)

func TryReadBotTokenFile() (string, error) {
	res, err := os.ReadFile("../bot.txt")

	if err != nil {
		return "", err
	}

	return string(res), nil
}
