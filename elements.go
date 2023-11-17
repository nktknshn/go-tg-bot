package tgbot

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

func AInputHandler[A any](handler func(string) any) *ElementInputHandler[A] {
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
	// String() string
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
	Handler func(string) any
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
	return "ElementCompleteMessage{}"
}

type ElementMessagePart struct {
	Text string
}

func (c *ElementMessagePart) elementKind() string {
	return KindElementMessagePart
}

func (c ElementMessagePart) String() string {
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

func (br *ElementButtonsRow[A]) Buttons() []ElementButton[A] {
	result := make([]ElementButton[A], 0)

	for idx, b := range br.Texts {
		idx := idx
		b := b
		result = append(result, ElementButton[A]{
			Text:    b,
			Action:  b,
			NextRow: false,
			OnClick: func() A {
				return br.OnClick(idx, b)
			},
		})
	}

	return result
}

func (c *ElementButtonsRow[A]) elementKind() string {
	return KindElementButtonsRow
}

func (c ElementButtonsRow[A]) String() string {
	return fmt.Sprintf("ElementButtonsRow{Texts=%v}", c.Texts)
}

type ElementBottomButton struct {
	Text  string
	Texts []string
	Hide  bool
}

func (c *ElementBottomButton) elementKind() string {
	return KindElementBottomButton
}

func (c ElementBottomButton) String() string {
	return fmt.Sprintf("ElementBottomButton{Text=%s, Texts=%s}", c.Text, c.Texts)
}

type ElementUserMessage struct {
	MessageID int
}

func (c *ElementUserMessage) elementKind() string {
	return KindElementUserMessage
}

func (c ElementUserMessage) String() string {
	return fmt.Sprintf("ElementUserMessage{MessageId=%d}", c.MessageID)
}

func (c *ElementUserMessage) Equal(other BasicElement) bool {

	if other.elementKind() != KindElementUserMessage {
		return false
	}

	otherUserMessage := other.(*ElementUserMessage)

	return c.MessageID == otherUserMessage.MessageID
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

type Elements []Element

func (es Elements) String() string {

	if len(es) == 0 {
		return "[]"
	}

	result := fmt.Sprintf("[%v", es[0])

	for _, e := range es[1:] {
		result = fmt.Sprintf("%v, %v", result, e)
	}

	return result + "]"
}

func (es Elements) ZapField(key string) zap.Field {
	return zap.Array(key, zapcore.ArrayMarshalerFunc(func(ae zapcore.ArrayEncoder) error {
		for _, e := range es {
			ae.AppendObject(
				zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {
					oe.AddString("kind", e.elementKind())
					return nil
				}),
			)
		}
		return nil
	}))
}
