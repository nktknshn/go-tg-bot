package btest

import (
	"testing"

	"github.com/BooleanCat/go-functional/iter"
	"github.com/nktknshn/go-tg-bot/emulator"
	"github.com/nktknshn/go-tg-bot/emulator/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertDisplayedMessages(t *testing.T, user *emulator.FakeBotUser, expectedMessages []helpers.MessageSimple) {
	messages := user.DisplayedMessages()

	if len(messages) != len(expectedMessages) {
		t.Errorf("expected %d messages, got %d", len(expectedMessages), len(messages))
		return
	}

	simpls := iter.Map(
		iter.Lift(messages),
		helpers.MessageAsSimple,
	).Collect()

	assert.ElementsMatch(t, expectedMessages, simpls)
}
