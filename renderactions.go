package tgbot

import (
	"context"
	"fmt"
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
	String() string
}

type RenderActionKeep struct {
	RenderedElement RenderedElement
	NewElement      OutcomingMessage
}

func (a *RenderActionKeep) RenderActionKind() string {
	return KindRenderActionKeep
}

func (a RenderActionKeep) String() string {
	return fmt.Sprintf("RenderActionKeep{RenderedElement: %v, NewElement: %v}", a.RenderedElement, a.NewElement)
}

type RenderActionReplace struct {
	RenderedElement RenderedElement
	NewElement      OutcomingMessage
}

func (a *RenderActionReplace) RenderActionKind() string {
	return KindRenderActionReplace
}

func (a RenderActionReplace) String() string {
	return fmt.Sprintf("RenderActionReplace{RenderedElement: %v, NewElement: %v}", a.RenderedElement, a.NewElement)
}

type RenderActionRemove struct {
	RenderedElement RenderedElement
}

func (a *RenderActionRemove) RenderActionKind() string {
	return KindRenderActionRemove
}

func (a RenderActionRemove) String() string {
	return fmt.Sprintf("RenderActionRemove{RenderedElement: %v}", a.RenderedElement)
}

type RenderActionCreate struct {
	NewElement OutcomingMessage
}

func (a *RenderActionCreate) RenderActionKind() string {
	return KindRenderActionCreate
}

func (a RenderActionCreate) String() string {
	return fmt.Sprintf("RenderActionCreate{NewElement: %v}", a.NewElement)
}

func areSame[A any](a RenderedElement, b OutcomingMessage) bool {

	if b == nil {
		return false
	}

	if a.renderedKind() == KindRenderedUserMessage && b.OutcomingKind() == KindOutcomingUserMessage {
		m := a.(*RenderedUserMessage)
		om := b.(*OutcomingUserMessage)

		return m.MessageId == om.ElementUserMessage.MessageID
	} else if a.renderedKind() == KindRenderedBotMessage && b.OutcomingKind() == KindOutcomingTextMessage {
		m := a.(*RenderedBotMessage[A])
		om := b.(*OutcomingTextMessage[A])

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
func GetRenderActions[A any](renderedElements []RenderedElement, nextElements []OutcomingMessage) []RenderActionType {

	logger := GetLogger()

	logger.Debug("GetRenderActions",
		zap.Any("renderedElements", len(renderedElements)),
		zap.Any("nextElements", len(nextElements)),
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

		// logger.Debug("GetRenderActions iteration",
		// 	zap.Any("r", r), zap.Any("n", n), zap.Any("idx", idx), zap.Any("len(result)", len(result)),
		// )

		if n == nil {
			// we are out of new elements to render so we can delete all remaining rendered elements
			result = slices.Delete(result, idx, idx+1)
			idx -= 1
			actions = append(actions, &RenderActionRemove{RenderedElement: r})
		} else if r == nil {
			// we are out of rendered elements so we can create all remaining new elements
			actions = append(actions, &RenderActionCreate{NewElement: n})
			continue
		} else if areSame[A](r, n) {
			actions = append(actions, &RenderActionKeep{RenderedElement: r, NewElement: n})
		} else if slices.IndexFunc(renderedElements, func(re RenderedElement) bool { return areSame[A](re, n) }) > idx {
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

const emptyString = "<empty>"

func GetOrText(text string, fallback string) string {
	if text == "" {
		return fallback
	}

	return text
}

func create[A any](ctx context.Context, renderer ChatRenderer, action *RenderActionCreate) (RenderedElement, error) {

	switch a := action.NewElement.(type) {
	case *OutcomingTextMessage[A]:

		message, err := renderer.Message(ctx, &ChatRendererMessageProps{
			Text:        GetOrText(a.Text, emptyString),
			ReplyMarkup: a.ReplyMarkup(),
		})

		if err != nil {
			return nil, err
		}

		return &RenderedBotMessage[A]{
			Message:              message,
			OutcomingTextMessage: a,
		}, nil

	case *OutcomingUserMessage:
		return &RenderedUserMessage{
			MessageId:            a.ElementUserMessage.MessageID,
			OutcomingUserMessage: a,
		}, nil
	// TODO
	// case *OutcomingFileMessage:
	// 	return renderer.File(a.ElementFile)
	// case *OutcomingPhotoGroupMessage:
	// 	return renderer.PhotoGroup(a.ElementPhotoGroup)
	default:
		globalLogger.Error("create: unsupported outcoming message type", zap.Any("a", a))
	}

	return nil, nil

}

// Takes actions and applies them to the renderer
func ExecuteRenderActions[A any](ctx context.Context, renderer ChatRenderer, actions []RenderActionType) ([]RenderedElement, error) {
	result := make([]RenderedElement, 0)
	actionsRemove := make([]RenderActionRemove, 0)
	actionsRemoveBot := make([]RenderActionRemove, 0)
	actionsRemoveUser := make([]RenderActionRemove, 0)

	actionsOther := make([]RenderActionType, 0)

	for _, action := range actions {
		switch a := action.(type) {
		case *RenderActionRemove:
			if a.RenderedElement.renderedKind() == KindRenderedBotMessage {
				actionsRemoveBot = append(actionsRemoveBot, *a)
			} else if a.RenderedElement.renderedKind() == KindRenderedUserMessage {
				actionsRemoveUser = append(actionsRemoveUser, *a)
			}
			actionsRemove = append(actionsRemove, *a)
		default:
			actionsOther = append(actionsOther, a)
		}
	}

	for _, a := range actionsRemoveUser {
		globalLogger.Debug("ExecuteRenderActions: removing rendered element", zap.Any("a", a))

		err := renderer.Delete(a.RenderedElement.ID())

		if err != nil {
			globalLogger.Error("Error removing rendered element", zap.Error(err))
		}
	}

	for _, action := range actionsOther {
		switch a := action.(type) {
		case *RenderActionCreate:
			globalLogger.Debug("ExecuteRenderActions: creating new element", zap.Any("a", a))

			rendereredMessage, err := create[A](ctx, renderer, a)

			if err != nil {
				return nil, err
			}

			result = append(result, rendereredMessage)
		case *RenderActionKeep:
			globalLogger.Debug("ExecuteRenderActions: keeping rendered element", zap.Any("a", a))

			if a.RenderedElement.renderedKind() == KindRenderedBotMessage && a.NewElement.OutcomingKind() == KindOutcomingTextMessage {
				rendereredMessage := &RenderedBotMessage[A]{
					OutcomingTextMessage: a.NewElement.(*OutcomingTextMessage[A]),
					Message:              a.RenderedElement.(*RenderedBotMessage[A]).Message,
				}

				result = append(result, rendereredMessage)
			}

		case *RenderActionReplace:
			if a.RenderedElement.renderedKind() == KindRenderedBotMessage && a.NewElement.OutcomingKind() == KindOutcomingTextMessage {

				globalLogger.Debug("ExecuteRenderActions: replacing rendered element", zap.Any("a", a.RenderedElement))

				outcoming := a.NewElement.(*OutcomingTextMessage[A])
				renderedElement := a.RenderedElement.(*RenderedBotMessage[A])

				globalLogger.Debug("ExecuteRenderActions: replacing rendered element",
					zap.Any("outcoming", outcoming), zap.Any("renderedElement", renderedElement),
				)

				message, err := renderer.Message(ctx, &ChatRendererMessageProps{
					Text:          GetOrText(outcoming.Text, emptyString),
					ReplyMarkup:   outcoming.ReplyMarkup(),
					TargetMessage: renderedElement.Message,
					RemoveTarget:  false,
				})

				if err != nil {
					return nil, err
				}

				rendereredMessage := &RenderedBotMessage[A]{
					OutcomingTextMessage: outcoming,
					Message:              message,
				}

				result = append(result, rendereredMessage)
			}
		}
	}

	for _, a := range actionsRemoveBot {
		globalLogger.Debug("ExecuteRenderActions: removing rendered element", zap.Any("a", a))

		err := renderer.Delete(a.RenderedElement.ID())

		if err != nil {
			globalLogger.Error("Error removing rendered element", zap.Error(err))
		}
	}

	// for _, action := range actionsRemove {
	// 	globalLogger.Debug("ExecuteRenderActions: removing rendered element", zap.Any("a", action))
	// 	err := renderer.Delete(action.RenderedElement.ID())

	// 	if err != nil {
	// 		globalLogger.Error("Error removing rendered element", zap.Error(err))
	// 	}
	// }

	return result, nil
}
