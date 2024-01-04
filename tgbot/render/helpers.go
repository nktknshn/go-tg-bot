package render

import "github.com/nktknshn/go-tg-bot/tgbot/outcoming"

// if a rendered element needs to be replaced with a new element
func isRenderedEqualOutcoming(a RenderedElement, b outcoming.OutcomingMessage) bool {

	if b == nil {
		return false
	}

	if a.RenderedKind() == KindRenderedUserMessage && b.OutcomingKind() == outcoming.KindOutcomingUserMessage {
		m := a.(*RenderedUserMessage)
		om := b.(*outcoming.OutcomingUserMessage)

		return m.MessageID == om.ElementUserMessage.MessageID
	} else if a.RenderedKind() == KindRenderedBotMessage && b.OutcomingKind() == outcoming.KindOutcomingTextMessage {
		m := a.(*RenderedBotMessage)
		om := b.(*outcoming.OutcomingTextMessage)

		return m.OutcomingTextMessage.Equal(om)

	} else if a.RenderedKind() == KindRenderedBotDocumentMessage && b.OutcomingKind() == outcoming.KindOutcomingFileMessage {
		m := a.(*RenderedBotDocumentMessage)
		om := b.(*outcoming.OutcomingFileMessage)

		return m.OutcomingFileMessage.ElementFile.FileId == om.ElementFile.FileId

	} else if a.RenderedKind() == KindRenderedPhotoGroup && b.OutcomingKind() == outcoming.KindOutcomingPhotoGroupMessage {
		// m := a.(*RenderedPhotoGroup)
		// om := b.(*OutcomingPhotoGroupMessage)
		// TODO: implement
		return false
	}

	return false
}
