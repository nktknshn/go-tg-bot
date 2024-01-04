package outcoming

// Elements -> Outcoming -> Rendered
// Messages that are going to be sent to user
const (
	KindOutcomingFileMessage       = "OutcomingFileMessage"
	KindOutcomingUserMessage       = "OutcomingUserMessage"
	KindOutcomingPhotoGroupMessage = "OutcomingPhotoGroupMessage"
	KindOutcomingTextMessage       = "OutcomingTextMessage"
)

// OutcomingMessage is a message that is going to be sent to user
// It's kind of compiled version of Elements
type OutcomingMessage interface {
	String() string
	OutcomingKind() string
	Equal(other OutcomingMessage) bool
}
