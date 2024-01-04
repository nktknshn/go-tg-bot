package outcoming

import (
	"fmt"

	"github.com/gotd/td/tg"
	"github.com/nktknshn/go-tg-bot/tgbot/component"
)

type OutcomingFileMessage struct {
	ElementFile component.ElementFile
	Message     *tg.Message
}

func (m OutcomingFileMessage) String() string {
	return fmt.Sprintf(
		"OutcomingFileMessage{ElementFile: %v, Message: %v}",
		m.ElementFile, m.Message,
	)
}

func (t *OutcomingFileMessage) OutcomingKind() string {
	return KindOutcomingFileMessage
}

func (t *OutcomingFileMessage) Equal(other OutcomingMessage) bool {
	if other.OutcomingKind() != KindOutcomingFileMessage {
		return false
	}

	otherFileMessage := other.(*OutcomingFileMessage)

	return t.ElementFile.FileId == otherFileMessage.ElementFile.FileId
}
