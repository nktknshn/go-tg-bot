package tgbot

import "fmt"

const (
	KindElementPhotoGroup      = "ElementPhotoGroup"
	KindElementFile            = "ElementFile"
	KindElementMessage         = "ElementMessage"
	KindElementMessagePart     = "ElementMessagePart"
	KindElementButton          = "ElementButton"
	KindElementBottomButton    = "KindElementBottomButton"
	KindElementComponent       = "ElementComponent"
	KindElementInputHandler    = "ElementInputHandler"
	KindElementUserMessage     = "ElementUserMessage"
	KindElementCompleteMessage = "ElementCompleteMessage"
	KindElementButtonsRow      = "ElementButtonsRow"
)

func MessageComplete() *ElementCompleteMessage {
	return &ElementCompleteMessage{}
}

func MessagePart(text string) *ElementMessagePart {
	return &ElementMessagePart{
		Text: text,
	}
}

func Message(text string) *ElementMessage {
	return &ElementMessage{
		Text: text,
	}
}

func Button[A any](text string, onClick func() A, action string, nextRow bool) *ElementButton[A] {
	return &ElementButton[A]{
		Text:    text,
		Action:  action,
		OnClick: onClick,
		NextRow: nextRow,
	}
}

func ButtonsRow[A any](texts []string, onClick func(int, string) A) *ElementButtonsRow[A] {
	return &ElementButtonsRow[A]{
		Texts:   texts,
		OnClick: onClick,
	}
}

func BottomButton(text string) *ElementBottomButton {
	return &ElementBottomButton{
		Text: text,
	}
}

func Component[A any](comp Comp[A]) *ElementComponent[A] {
	return &ElementComponent[A]{
		comp: comp,
	}
}

func AInputHandler[A any](handler func(string) A) *ElementInputHandler[A] {
	return &ElementInputHandler[A]{
		Handler: handler,
	}
}

// Element is button, message, handler etc...
type BasicElement interface {
	String() string
	elementKind() string
}

// Element is BasicElement or another component etc...
type Element interface {
	String() string
	elementKind() string
}

type ElementComponent[A any] struct {
	comp Comp[A]
}

func (c *ElementComponent[A]) elementKind() string {
	return KindElementComponent
}

func (c ElementComponent[A]) String() string {
	return "ElementComponent{...}"
}

type ElementInputHandler[A any] struct {
	Handler func(string) A
}

func (c *ElementInputHandler[A]) elementKind() string {
	return KindElementInputHandler
}

func (c ElementInputHandler[A]) String() string {
	return "ElementInputHandler{}"
}

type ElementPhotoGroup struct {
	photos []ElementFile
}

func (c *ElementPhotoGroup) elementKind() string {
	return KindElementPhotoGroup
}

func (c ElementPhotoGroup) String() string {
	return fmt.Sprintf("ElementPhotoGroup{photos=%v}", c.photos)
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

func (ef ElementFile) String() string {
	return fmt.Sprintf("ElementFile{FileId=%s}", ef.FileId)
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

func (c ElementMessage) String() string {
	return fmt.Sprintf("ElementMessage{Text=%s}", c.Text)
}

func (c *ElementMessage) Equal(other BasicElement) bool {
	if other.elementKind() != KindElementMessage {
		return false
	}

	otherMessage := other.(*ElementMessage)

	return c.Text == otherMessage.Text
}

type ElementCompleteMessage struct{}

func (c *ElementCompleteMessage) elementKind() string {
	return KindElementCompleteMessage
}

func (c ElementCompleteMessage) String() string {
	return fmt.Sprintf("ElementCompleteMessage{}")
}

type ElementMessagePart struct {
	Text string
}

func (c *ElementMessagePart) elementKind() string {
	return KindElementMessagePart
}

func (c *ElementMessagePart) String() string {
	return fmt.Sprintf("ElementMessagePart{Text=%s}", c.Text)
}

func (c *ElementMessagePart) Equal(other BasicElement) bool {
	if other.elementKind() != KindElementMessagePart {
		return false
	}

	otherMessagePart := other.(*ElementMessagePart)

	return c.Text == otherMessagePart.Text
}

type ElementButton[A any] struct {
	Text    string
	Action  string
	NextRow bool
	OnClick func() A
}

func (b *ElementButton[A]) CallbackData() string {
	if b.Action != "" {
		return b.Action
	}

	return b.Text
}

func (c *ElementButton[A]) elementKind() string {
	return KindElementButton
}

func (c ElementButton[A]) String() string {
	return fmt.Sprintf("ElementButton{Text=%s, Action=%s}", c.Text, c.Action)
}

func (c *ElementButton[A]) Equal(other BasicElement) bool {
	if other.elementKind() != KindElementButton {
		return false
	}

	otherButton := other.(*ElementButton[A])

	return c.Text == otherButton.Text && c.Action == otherButton.Action
}

type ElementButtonsRow[A any] struct {
	Texts   []string
	OnClick func(int, string) A
}

func (c *ElementButtonsRow[A]) elementKind() string {
	return KindElementButtonsRow
}

func (c ElementButtonsRow[A]) String() string {
	return fmt.Sprintf("ElementButtonsRow{Texts=%v}", c.Texts)
}

type ElementBottomButton struct {
	Text string
}

func (c *ElementBottomButton) elementKind() string {
	return KindElementBottomButton
}

func (c ElementBottomButton) String() string {
	return fmt.Sprintf("ElementBottomButton{Text=%s}", c.Text)
}

type ElementUserMessage struct {
	MessageId int
}

func (c *ElementUserMessage) elementKind() string {
	return KindElementUserMessage
}

func (c ElementUserMessage) String() string {
	return fmt.Sprintf("ElementUserMessage{MessageId=%d}", c.MessageId)
}

func (c *ElementUserMessage) Equal(other BasicElement) bool {

	if other.elementKind() != KindElementUserMessage {
		return false
	}

	otherUserMessage := other.(*ElementUserMessage)

	return c.MessageId == otherUserMessage.MessageId
}
