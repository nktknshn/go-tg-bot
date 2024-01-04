package outcoming

import (
	"fmt"

	"github.com/nktknshn/go-tg-bot/tgbot/component"
)

type OutcomingUserMessage struct {
	ElementUserMessage component.ElementUserMessage
}

func (m OutcomingUserMessage) String() string {
	return fmt.Sprintf(
		"OutcomingUserMessage{ElementUserMessage: %v}",
		m.ElementUserMessage,
	)
}

func (t *OutcomingUserMessage) OutcomingKind() string {
	return KindOutcomingUserMessage
}

func (t *OutcomingUserMessage) Equal(other OutcomingMessage) bool {
	if other.OutcomingKind() != KindOutcomingUserMessage {
		return false
	}

	otherUserMessage := other.(*OutcomingUserMessage)

	return t.ElementUserMessage.MessageID == otherUserMessage.ElementUserMessage.MessageID
}
