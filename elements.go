package tgbot

const (
	KindElementPhotoGroup   = "ElementPhotoGroup"
	KindElementFile         = "ElementFile"
	KindElementMessage      = "ElementMessage"
	KindElementMessagePart  = "ElementMessagePart"
	KindElementButton       = "ElementButton"
	KindElementComponent    = "ElementComponent"
	KindElementInputHandler = "ElementInputHandler"
	KindElementUserMessage  = "ElementUserMessage"
)

func EndMessage() *ElementMessagePart {
	return &ElementMessagePart{
		text: "",
	}
}

func MessagePart(text string) *ElementMessagePart {
	return &ElementMessagePart{
		text: text,
	}
}

func Message(text string) *ElementMessage {
	return &ElementMessage{
		Text: text,
	}
}

func Button[A any](text string, onClick func() A) *ElementButton[A] {
	return &ElementButton[A]{
		Text:    text,
		OnClick: onClick,
	}
}

func Component[A any](comp Comp[A]) *ElementComponent[A] {
	return &ElementComponent[A]{
		comp: comp,
	}
}

func AInputHandler[A any](handler InputHandler[A]) *ElementInputHandler[A] {
	return &ElementInputHandler[A]{
		Handler: handler,
	}
}

// Element is button, message, handler etc...
type BasicElement interface {
	elementKind() string
	// Equal(other BasicElement) bool
}

// Element is BasicElement or another component etc...
type Element interface {
	elementKind() string
}

// func (e Element) String() string {
// 	return e.elementKind()
// }

type ElementComponent[A any] struct {
	comp Comp[A]
}

func (c *ElementComponent[A]) elementKind() string {
	return KindElementComponent
}

type ElementInputHandler[A any] struct {
	Handler InputHandler[A]
}

func (c *ElementInputHandler[A]) elementKind() string {
	return KindElementInputHandler
}

type ElementPhotoGroup struct {
	photos []ElementFile
}

func (c *ElementPhotoGroup) elementKind() string {
	return KindElementPhotoGroup
}

func (c *ElementPhotoGroup) Equal(other BasicElement) bool {
	if other.elementKind() != KindElementPhotoGroup {
		return false
	}

	otherPhotoGroup := other.(*ElementPhotoGroup)

	if len(c.photos) != len(otherPhotoGroup.photos) {
		return false
	}

	for i, photo := range c.photos {
		if !photo.Equal(otherPhotoGroup.photos[i]) {
			return false
		}
	}

	return true
}

type ElementFile struct {
	FileId string
}

func (c ElementFile) elementKind() string {
	return KindElementFile
}

func (c *ElementFile) Equal(other BasicElement) bool {
	if other.elementKind() != KindElementFile {
		return false
	}

	otherFile := other.(*ElementFile)

	return c.FileId == otherFile.FileId
}

type ElementMessage struct {
	Text string
}

func (c *ElementMessage) elementKind() string {
	return KindElementMessage
}

type ElementMessagePart struct {
	text string
}

func (c *ElementMessagePart) elementKind() string {
	return KindElementMessagePart
}

type ElementButton[A any] struct {
	Text    string
	Action  string
	OnClick func() A
}

func (c *ElementButton[A]) elementKind() string {
	return KindElementButton
}

type ElementUserMessage struct {
	MessageId int
}

func (c *ElementUserMessage) elementKind() string {
	return KindElementUserMessage
}
