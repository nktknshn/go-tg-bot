package outcoming

import (
	"fmt"

	"github.com/nktknshn/go-tg-bot/tgbot/component"
)

type OutcomingPhotoGroupMessage struct {
	ElementPhotoGroup component.ElementPhotoGroup
}

func (m OutcomingPhotoGroupMessage) String() string {
	return fmt.Sprintf(
		"OutcomingPhotoGroupMessage{ElementPhotoGroup: %v}",
		m.ElementPhotoGroup,
	)
}

func (t *OutcomingPhotoGroupMessage) OutcomingKind() string {
	return KindOutcomingPhotoGroupMessage
}

func (t *OutcomingPhotoGroupMessage) Equal(other OutcomingMessage) bool {
	if other.OutcomingKind() != KindOutcomingPhotoGroupMessage {
		return false
	}

	otherPhotoGroupMessage := other.(*OutcomingPhotoGroupMessage)

	return otherPhotoGroupMessage.ElementPhotoGroup.Equal(&t.ElementPhotoGroup)
}
