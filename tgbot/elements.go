package tgbot

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

func newMessageComplete() *elementCompleteMessage {
	return &elementCompleteMessage{}
}

func newMessagePart(text string) *elementMessagePart {
	return &elementMessagePart{
		Text: text,
	}
}

func newMessage(text string) *elementMessage {
	return &elementMessage{
		Text: text,
	}
}

func newButton(text string, onClick func() any, action string, nextRow bool, noCallback bool) *elementButton {
	return &elementButton{
		Text:       text,
		Action:     action,
		OnClick:    onClick,
		NextRow:    nextRow,
		NoCallback: noCallback,
	}
}

func newButtonsRow(texts []string, onClick func(int, string) any) *elementButtonsRow {
	return &elementButtonsRow{
		Texts:   texts,
		OnClick: onClick,
	}
}

func newBottomButton(text string) *elementBottomButton {
	return &elementBottomButton{
		Text: text,
	}
}

func newComponent(comp Comp) *elementComponent {
	return &elementComponent{
		comp: comp,
	}
}

func newInputHandler(handler func(string) any) *elementInputHandler {
	return &elementInputHandler{
		Handler: handler,
	}
}

// Element is button, message, handler etc...
type basicElement interface {
	String() string
	elementKind() string
}

// anyElement is BasicElement or another component etc...
type anyElement interface {
	// String() string
	elementKind() string
}

type elementComponent struct {
	comp Comp
}

func (c *elementComponent) elementKind() string {
	return kindElementComponent
}

func (c elementComponent) String() string {
	return fmt.Sprintf("ElementComponent{comp=%v}", reflectCompId(c.comp))
}

type elementInputHandler struct {
	Handler func(string) any
}

func (c *elementInputHandler) elementKind() string {
	return kindElementInputHandler
}

func (c elementInputHandler) String() string {
	return "ElementInputHandler{}"
}

type elementPhotoGroup struct {
	photos []elementFile
}

func (c *elementPhotoGroup) elementKind() string {
	return kindElementPhotoGroup
}

func (c elementPhotoGroup) String() string {
	return fmt.Sprintf("ElementPhotoGroup{photos=%v}", c.photos)
}

func (c *elementPhotoGroup) Equal(other basicElement) bool {
	if other.elementKind() != kindElementPhotoGroup {
		return false
	}

	otherPhotoGroup := other.(*elementPhotoGroup)

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

type elementMessage struct {
	Text string
}

func (c *elementMessage) elementKind() string {
	return kindElementMessage
}

func (c elementMessage) String() string {
	return fmt.Sprintf("ElementMessage{Text=%s}", c.Text)
}

func (c *elementMessage) Equal(other basicElement) bool {
	if other.elementKind() != kindElementMessage {
		return false
	}

	otherMessage := other.(*elementMessage)

	return c.Text == otherMessage.Text
}

type elementCompleteMessage struct{}

func (c *elementCompleteMessage) elementKind() string {
	return kindElementCompleteMessage
}

func (c elementCompleteMessage) String() string {
	return "ElementCompleteMessage{}"
}

type elementMessagePart struct {
	Text string
}

func (c *elementMessagePart) elementKind() string {
	return kindElementMessagePart
}

func (c elementMessagePart) String() string {
	return fmt.Sprintf("ElementMessagePart{Text=%s}", c.Text)
}

func (c *elementMessagePart) Equal(other basicElement) bool {
	if other.elementKind() != kindElementMessagePart {
		return false
	}

	otherMessagePart := other.(*elementMessagePart)

	return c.Text == otherMessagePart.Text
}

type elementButton struct {
	Text       string
	Action     string
	NextRow    bool
	NoCallback bool
	OnClick    func() any
}

func (b *elementButton) CallbackData() string {
	if b.Action != "" {
		return b.Action
	}

	return b.Text
}

func (c *elementButton) elementKind() string {
	return kindElementButton
}

func (c elementButton) String() string {
	return fmt.Sprintf("ElementButton{Text=%s, Action=%s}", c.Text, c.Action)
}

func (c *elementButton) Equal(other basicElement) bool {
	if other.elementKind() != kindElementButton {
		return false
	}

	otherButton := other.(*elementButton)

	return c.Text == otherButton.Text && c.Action == otherButton.Action
}

type elementButtonsRow struct {
	Texts   []string
	OnClick func(int, string) any
}

func (br *elementButtonsRow) Buttons() []elementButton {
	result := make([]elementButton, 0)

	for idx, b := range br.Texts {
		idx := idx
		b := b
		result = append(result, elementButton{
			Text:    b,
			Action:  b,
			NextRow: false,
			OnClick: func() any {
				return br.OnClick(idx, b)
			},
		})
	}

	return result
}

func (c *elementButtonsRow) elementKind() string {
	return kindElementButtonsRow
}

func (c elementButtonsRow) String() string {
	return fmt.Sprintf("ElementButtonsRow{Texts=%v}", c.Texts)
}

type elementBottomButton struct {
	Text  string
	Texts []string
	Hide  bool
}

func (c *elementBottomButton) elementKind() string {
	return kindElementBottomButton
}

func (c elementBottomButton) String() string {
	return fmt.Sprintf("ElementBottomButton{Text=%s, Texts=%s}", c.Text, c.Texts)
}

type elementUserMessage struct {
	MessageID int
}

func (c *elementUserMessage) elementKind() string {
	return kindElementUserMessage
}

func (c elementUserMessage) String() string {
	return fmt.Sprintf("ElementUserMessage{MessageId=%d}", c.MessageID)
}

func (c *elementUserMessage) Equal(other basicElement) bool {

	if other.elementKind() != kindElementUserMessage {
		return false
	}

	otherUserMessage := other.(*elementUserMessage)

	return c.MessageID == otherUserMessage.MessageID
}

type elementFile struct {
	FileId string
}

func (ef elementFile) String() string {
	return fmt.Sprintf("ElementFile{FileId=%s}", ef.FileId)
}

func (c elementFile) elementKind() string {
	return kindElementFile
}

func (c *elementFile) Equal(other basicElement) bool {
	if other.elementKind() != kindElementFile {
		return false
	}

	otherFile := other.(*elementFile)

	return c.FileId == otherFile.FileId
}

type elementsList []anyElement

func (es elementsList) String() string {

	if len(es) == 0 {
		return "[]"
	}

	result := fmt.Sprintf("[%v", es[0])

	for _, e := range es[1:] {
		result = fmt.Sprintf("%v, %v", result, e)
	}

	return result + "]"
}

func (es elementsList) ZapField(key string) zap.Field {
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
