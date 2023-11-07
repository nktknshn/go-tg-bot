package tgbot

import (
	"slices"

	"go.uber.org/zap"
)

const (
	KindRenderActionKeep    = "RenderActionKeep"
	KindRenderActionReplace = "RenderActionReplace"
	KindRenderActionRemove  = "RenderActionRemove"
	KindRenderActionCreate  = "RenderActionCreate"
)

type RenderActionType interface {
	RenderActionKind() string
}

type RenderActionKeep struct {
	RenderedElement RenderedElement
	NewElement      OutcomingMessage
}

func (a *RenderActionKeep) RenderActionKind() string {
	return KindRenderActionKeep
}

type RenderActionReplace struct {
	RenderedElement RenderedElement
	NewElement      OutcomingMessage
}

func (a *RenderActionReplace) RenderActionKind() string {
	return KindRenderActionReplace
}

type RenderActionRemove struct {
	RenderedElement RenderedElement
}

func (a *RenderActionRemove) RenderActionKind() string {
	return KindRenderActionRemove
}

type RenderActionCreate struct {
	NewElement OutcomingMessage
}

func (a *RenderActionCreate) RenderActionKind() string {
	return KindRenderActionCreate
}

func areSame(a RenderedElement, b OutcomingMessage) bool {

	if b == nil {
		return false
	}

	if a.renderedKind() == KindRenderedUserMessage && b.OutcomingKind() == KindOutcomingUserMessage {
		m := a.(*RenderedUserMessage)
		om := b.(*OutcomingUserMessage)

		return m.MessageId == om.ElementUserMessage.MessageId
	} else if a.renderedKind() == KindRenderedBotMessage && b.OutcomingKind() == KindOutcomingTextMessage {
		m := a.(*RenderedBotMessage)
		om := b.(*OutcomingTextMessage[any])

		return m.OutcomingTextMessage.Equal(om)

	} else if a.renderedKind() == KindRenderedBotDocumentMessage && b.OutcomingKind() == KindOutcomingFileMessage {
		m := a.(*RenderedBotDocumentMessage)
		om := b.(*OutcomingFileMessage)

		return m.OutcomingFileMessage.ElementFile.FileId == om.ElementFile.FileId

	} else if a.renderedKind() == KindRenderedPhotoGroup && b.OutcomingKind() == KindOutcomingPhotoGroupMessage {
		// m := a.(*RenderedPhotoGroup)
		// om := b.(*OutcomingPhotoGroupMessage)
		// TODO: implement
		return false
	}

	return false
}

/*
Rules are:
1. never talk about the fight club
*/
func GetRenderActions(renderedElements []RenderedElement, nextElements []OutcomingMessage) []RenderActionType {

	logger := GetLogger()

	logger.Debug("GetRenderActions",
		zap.Any("renderedElements", renderedElements),
		zap.Any("nextElements", nextElements),
	)

	actions := make([]RenderActionType, 0)

	// renderedElements = append(make([]RenderedElement, 0), renderedElements...)
	// nextElements = append(make([]OutcomingMessage, 0), nextElements...)

	result := append(make([]RenderedElement, 0), renderedElements...)
	idx := 0

	for {

		outOfRenderedElements := (idx > len(result)-1) || idx < 0

		// do while we have either rendered elements or new elements
		// if we have no rendered elements and no new elements we are done
		if outOfRenderedElements && len(nextElements) == 0 {
			break
		}

		var r RenderedElement
		var n OutcomingMessage

		if len(result) > idx {
			r = result[idx]
		} else {
			r = nil
		}

		if len(nextElements) > 0 {
			n = nextElements[0]
			nextElements = nextElements[1:]
		} else {
			n = nil
		}

		logger.Debug("GetRenderActions iteration",
			zap.Any("r", r), zap.Any("n", n), zap.Any("idx", idx), zap.Any("len(result)", len(result)),
		)

		if n == nil {
			// we are out of new elements to render so we can delete all remaining rendered elements
			result = slices.Delete(result, idx, idx+1)
			idx -= 1
			actions = append(actions, &RenderActionRemove{RenderedElement: r})
		} else if r == nil {
			// we are out of rendered elements so we can create all remaining new elements
			actions = append(actions, &RenderActionCreate{NewElement: n})
			continue
		} else if areSame(r, n) {
			actions = append(actions, &RenderActionKeep{RenderedElement: r, NewElement: n})
		} else if slices.IndexFunc(renderedElements, func(re RenderedElement) bool { return areSame(re, n) }) > idx {
			// if we have the next outcoming element rendered somewhere else ahead of current rendered element
			// we can delete current rendered element
			result = slices.Delete(result, idx, idx+1)
			nextElements = append([]OutcomingMessage{n}, nextElements...)
			idx -= 1
			actions = append(actions, &RenderActionRemove{RenderedElement: r})
		} else {
			if r.canReplace(n) {
				actions = append(actions, &RenderActionReplace{RenderedElement: r, NewElement: n})
			} else {
				result = slices.Delete(result, idx, idx+1)
				nextElements = append([]OutcomingMessage{n}, nextElements...)
				idx -= 1
				actions = append(actions, &RenderActionRemove{RenderedElement: r})
			}
		}

		idx += 1

	}

	return actions
}
