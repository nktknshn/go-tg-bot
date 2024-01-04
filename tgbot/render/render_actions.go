package render

import (
	"fmt"

	"github.com/nktknshn/go-tg-bot/tgbot/outcoming"
)

const (
	kindRenderActionKeep    = "RenderActionKeep"
	kindRenderActionReplace = "RenderActionReplace"
	kindRenderActionRemove  = "RenderActionRemove"
	kindRenderActionCreate  = "RenderActionCreate"
)

type RenderAction interface {
	RenderActionKind() string
	String() string
}

type renderActionKeep struct {
	RenderedElement RenderedElement
	NewElement      outcoming.OutcomingMessage
}

func (a *renderActionKeep) RenderActionKind() string {
	return kindRenderActionKeep
}

func (a renderActionKeep) String() string {
	return fmt.Sprintf("RenderActionKeep{RenderedElement: %v, NewElement: %v}", a.RenderedElement, a.NewElement)
}

type renderActionReplace struct {
	RenderedElement RenderedElement
	NewElement      outcoming.OutcomingMessage
}

func (a *renderActionReplace) RenderActionKind() string {
	return kindRenderActionReplace
}

func (a renderActionReplace) String() string {
	return fmt.Sprintf("RenderActionReplace{RenderedElement: %v, NewElement: %v}", a.RenderedElement, a.NewElement)
}

type renderActionRemove struct {
	RenderedElement RenderedElement
}

func (a *renderActionRemove) RenderActionKind() string {
	return kindRenderActionRemove
}

func (a renderActionRemove) String() string {
	return fmt.Sprintf("RenderActionRemove{RenderedElement: %v}", a.RenderedElement)
}

type renderActionCreate struct {
	NewElement outcoming.OutcomingMessage
}

func (a *renderActionCreate) RenderActionKind() string {
	return kindRenderActionCreate
}

func (a renderActionCreate) String() string {
	return fmt.Sprintf("RenderActionCreate{NewElement: %v}", a.NewElement)
}
