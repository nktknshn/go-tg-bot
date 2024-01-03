package component

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	kindElementPhotoGroup      = "ElementPhotoGroup"
	kindElementFile            = "ElementFile"
	kindElementMessage         = "ElementMessage"
	kindElementMessagePart     = "ElementMessagePart"
	kindElementButton          = "ElementButton"
	kindElementBottomButton    = "ElementBottomButton"
	kindElementComponent       = "ElementComponent"
	kindElementInputHandler    = "ElementInputHandler"
	kindElementUserMessage     = "ElementUserMessage"
	kindElementCompleteMessage = "ElementCompleteMessage"
	kindElementButtonsRow      = "ElementButtonsRow"
)

func newMessageComplete() *ElementCompleteMessage {
	return &ElementCompleteMessage{}
}

func newMessagePart(text string) *ElementMessagePart {
	return &ElementMessagePart{
		Text: text,
	}
}

func newMessage(text string) *ElementMessage {
	return &ElementMessage{
		Text: text,
	}
}

func NewButton(text string, onClick func() any, action string, nextRow bool, noCallback bool) *ElementButton {
	return &ElementButton{
		Text:       text,
		Data:       action,
		OnClick:    onClick,
		NextRow:    nextRow,
		NoCallback: noCallback,
	}
}

func newButtonsRow(texts []string, onClick func(int, string) any) *ElementButtonsRow {
	return &ElementButtonsRow{
		Texts:   texts,
		OnClick: onClick,
	}
}

func newBottomButton(text string) *ElementBottomButton {
	return &ElementBottomButton{
		Text: text,
	}
}

func newComponent(comp Comp) *ElementComponent {
	return &ElementComponent{
		comp: comp,
	}
}

func newInputHandler(handler func(string) any) *ElementInputHandler {
	return &ElementInputHandler{
		Handler: handler,
	}
}

// Element is button, message, handler etc...
type BasicElement interface {
	String() string
	elementKind() string
}

// AnyElement is BasicElement or another component etc...
type AnyElement interface {
	// String() string
	elementKind() string
}

type ElementComponent struct {
	comp Comp
}

func (c *ElementComponent) elementKind() string {
	return kindElementComponent
}

func (c ElementComponent) String() string {
	return fmt.Sprintf("ElementComponent{comp=%v}", reflectCompId(c.comp))
}

type ElementInputHandler struct {
	Handler func(string) any
}

func (c *ElementInputHandler) elementKind() string {
	return kindElementInputHandler
}

func (c ElementInputHandler) String() string {
	return "ElementInputHandler{}"
}

type ElementPhotoGroup struct {
	photos []ElementFile
}

func (c *ElementPhotoGroup) elementKind() string {
	return kindElementPhotoGroup
}

func (c ElementPhotoGroup) String() string {
	return fmt.Sprintf("ElementPhotoGroup{photos=%v}", c.photos)
}

func (c *ElementPhotoGroup) Equal(other BasicElement) bool {
	if other.elementKind() != kindElementPhotoGroup {
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
	return kindElementMessage
}

func (c ElementMessage) String() string {
	return fmt.Sprintf("ElementMessage{Text=%s}", c.Text)
}

func (c *ElementMessage) Equal(other BasicElement) bool {
	if other.elementKind() != kindElementMessage {
		return false
	}

	otherMessage := other.(*ElementMessage)

	return c.Text == otherMessage.Text
}

type ElementCompleteMessage struct{}

func (c *ElementCompleteMessage) elementKind() string {
	return kindElementCompleteMessage
}

func (c ElementCompleteMessage) String() string {
	return "ElementCompleteMessage{}"
}

type ElementMessagePart struct {
	Text string
}

func (c *ElementMessagePart) elementKind() string {
	return kindElementMessagePart
}

func (c ElementMessagePart) String() string {
	return fmt.Sprintf("ElementMessagePart{Text=%s}", c.Text)
}

func (c *ElementMessagePart) Equal(other BasicElement) bool {
	if other.elementKind() != kindElementMessagePart {
		return false
	}

	otherMessagePart := other.(*ElementMessagePart)

	return c.Text == otherMessagePart.Text
}

type ElementButton struct {
	Text       string
	Data       string
	NextRow    bool
	NoCallback bool
	OnClick    func() any
}

func (b *ElementButton) CallbackData() string {
	if b.Data != "" {
		return b.Data
	}

	return b.Text
}

func (c *ElementButton) elementKind() string {
	return kindElementButton
}

func (c ElementButton) String() string {
	return fmt.Sprintf("ElementButton{Text=%s, Action=%s}", c.Text, c.Data)
}

func (c *ElementButton) Equal(other BasicElement) bool {
	if other.elementKind() != kindElementButton {
		return false
	}

	otherButton := other.(*ElementButton)

	return c.Text == otherButton.Text && c.Data == otherButton.Data
}

type ElementButtonsRow struct {
	Texts   []string
	OnClick func(int, string) any
}

func (br *ElementButtonsRow) Buttons() []ElementButton {
	result := make([]ElementButton, 0)

	for idx, b := range br.Texts {
		idx := idx
		b := b
		result = append(result, ElementButton{
			Text:    b,
			Data:    b,
			NextRow: false,
			OnClick: func() any {
				return br.OnClick(idx, b)
			},
		})
	}

	return result
}

func (c *ElementButtonsRow) elementKind() string {
	return kindElementButtonsRow
}

func (c ElementButtonsRow) String() string {
	return fmt.Sprintf("ElementButtonsRow{Texts=%v}", c.Texts)
}

type ElementBottomButton struct {
	Text  string
	Texts []string
	Hide  bool
}

func (c *ElementBottomButton) elementKind() string {
	return kindElementBottomButton
}

func (c ElementBottomButton) String() string {
	return fmt.Sprintf("ElementBottomButton{Text=%s, Texts=%s}", c.Text, c.Texts)
}

type ElementUserMessage struct {
	MessageID int
}

func (c *ElementUserMessage) elementKind() string {
	return kindElementUserMessage
}

func (c ElementUserMessage) String() string {
	return fmt.Sprintf("ElementUserMessage{MessageId=%d}", c.MessageID)
}

func (c *ElementUserMessage) Equal(other BasicElement) bool {

	if other.elementKind() != kindElementUserMessage {
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
	return kindElementFile
}

func (c *ElementFile) Equal(other BasicElement) bool {
	if other.elementKind() != kindElementFile {
		return false
	}

	otherFile := other.(*ElementFile)

	return c.FileId == otherFile.FileId
}

type ElementsList []AnyElement

func (es ElementsList) String() string {

	if len(es) == 0 {
		return "[]"
	}

	result := fmt.Sprintf("[%v", es[0])

	for _, e := range es[1:] {
		result = fmt.Sprintf("%v, %v", result, e)
	}

	return result + "]"
}

func (es ElementsList) ZapField(key string) zap.Field {
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
