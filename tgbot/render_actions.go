package tgbot

import (
	"context"
	"fmt"
	"slices"

	"go.uber.org/zap"
)

const (
	kindRenderActionKeep    = "RenderActionKeep"
	kindRenderActionReplace = "RenderActionReplace"
	kindRenderActionRemove  = "RenderActionRemove"
	kindRenderActionCreate  = "RenderActionCreate"
)

type renderActionType interface {
	RenderActionKind() string
	String() string
}

type renderActionKeep struct {
	RenderedElement RenderedElement
	NewElement      outcomingMessage
}

func (a *renderActionKeep) RenderActionKind() string {
	return kindRenderActionKeep
}

func (a renderActionKeep) String() string {
	return fmt.Sprintf("RenderActionKeep{RenderedElement: %v, NewElement: %v}", a.RenderedElement, a.NewElement)
}

type renderActionReplace struct {
	RenderedElement RenderedElement
	NewElement      outcomingMessage
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
	NewElement outcomingMessage
}

func (a *renderActionCreate) RenderActionKind() string {
	return kindRenderActionCreate
}

func (a renderActionCreate) String() string {
	return fmt.Sprintf("RenderActionCreate{NewElement: %v}", a.NewElement)
}

func areSame(a RenderedElement, b outcomingMessage) bool {

	if b == nil {
		return false
	}

	if a.renderedKind() == kindRenderedUserMessage && b.OutcomingKind() == kindOutcomingUserMessage {
		m := a.(*renderedUserMessage)
		om := b.(*outcomingUserMessage)

		return m.MessageID == om.ElementUserMessage.MessageID
	} else if a.renderedKind() == kindRenderedBotMessage && b.OutcomingKind() == kindOutcomingTextMessage {
		m := a.(*renderedBotMessage)
		om := b.(*outcomingTextMessage)

		return m.OutcomingTextMessage.Equal(om)

	} else if a.renderedKind() == kindRenderedBotDocumentMessage && b.OutcomingKind() == kindOutcomingFileMessage {
		m := a.(*renderedBotDocumentMessage)
		om := b.(*outcomingFileMessage)

		return m.OutcomingFileMessage.ElementFile.FileId == om.ElementFile.FileId

	} else if a.renderedKind() == kindRenderedPhotoGroup && b.OutcomingKind() == kindOutcomingPhotoGroupMessage {
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
func getRenderActions(renderedElements []RenderedElement, nextElements []outcomingMessage, logger *zap.Logger) []renderActionType {

	logger.Debug("GetRenderActions",
		zap.Any("renderedElements", len(renderedElements)),
		zap.Any("nextElements", len(nextElements)),
	)

	actions := make([]renderActionType, 0)

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
		var n outcomingMessage

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
			actions = append(actions, &renderActionRemove{RenderedElement: r})
		} else if r == nil {
			// we are out of rendered elements so we can create all remaining new elements
			actions = append(actions, &renderActionCreate{NewElement: n})
			continue
		} else if areSame(r, n) {
			actions = append(actions, &renderActionKeep{RenderedElement: r, NewElement: n})
		} else if slices.IndexFunc(renderedElements, func(re RenderedElement) bool { return areSame(re, n) }) > idx {
			// if we have the next outcoming element rendered somewhere else ahead of current rendered element
			// we can delete current rendered element
			result = slices.Delete(result, idx, idx+1)
			nextElements = append([]outcomingMessage{n}, nextElements...)
			idx -= 1
			actions = append(actions, &renderActionRemove{RenderedElement: r})
		} else {
			if r.canReplace(n) {
				actions = append(actions, &renderActionReplace{RenderedElement: r, NewElement: n})
			} else {
				result = slices.Delete(result, idx, idx+1)
				nextElements = append([]outcomingMessage{n}, nextElements...)
				idx -= 1
				actions = append(actions, &renderActionRemove{RenderedElement: r})
			}
		}

		idx += 1

	}

	return actions
}

const emptyString = "<empty>"

func getOrText(text string, fallback string) string {
	if text == "" {
		return fallback
	}

	return text
}

func create(ctx context.Context, renderer ChatRenderer, action *renderActionCreate, logger *zap.Logger) (RenderedElement, error) {

	switch a := action.NewElement.(type) {
	case *outcomingTextMessage:

		message, err := renderer.Message(ctx, &ChatRendererMessageProps{
			Text:        getOrText(a.Text, emptyString),
			ReplyMarkup: a.ReplyMarkup(),
		})

		if err != nil {
			return nil, err
		}

		return &renderedBotMessage{
			Message:              message,
			OutcomingTextMessage: a,
		}, nil

	case *outcomingUserMessage:
		return &renderedUserMessage{
			MessageID:            a.ElementUserMessage.MessageID,
			OutcomingUserMessage: a,
		}, nil
	// TODO
	// case *OutcomingFileMessage:
	// 	return renderer.File(a.ElementFile)
	// case *OutcomingPhotoGroupMessage:
	// 	return renderer.PhotoGroup(a.ElementPhotoGroup)
	default:
		logger.Error("create: unsupported outcoming message type", zap.Any("a", a))
	}

	return nil, nil

}

// Takes actions and applies them to the renderer
func ExecuteRenderActions(ctx context.Context, renderer ChatRenderer, actions []renderActionType, logger *zap.Logger) ([]RenderedElement, error) {
	result := make([]RenderedElement, 0)
	actionsRemove := make([]renderActionRemove, 0)
	actionsRemoveBot := make([]renderActionRemove, 0)
	actionsRemoveUser := make([]renderActionRemove, 0)

	actionsOther := make([]renderActionType, 0)

	for _, action := range actions {
		switch a := action.(type) {
		case *renderActionRemove:
			if a.RenderedElement.renderedKind() == kindRenderedBotMessage {
				actionsRemoveBot = append(actionsRemoveBot, *a)
			} else if a.RenderedElement.renderedKind() == kindRenderedUserMessage {
				actionsRemoveUser = append(actionsRemoveUser, *a)
			}
			actionsRemove = append(actionsRemove, *a)
		default:
			actionsOther = append(actionsOther, a)
		}
	}

	for _, a := range actionsRemoveUser {
		logger.Debug("ExecuteRenderActions: removing rendered element", zap.Any("a", a))

		err := renderer.Delete(a.RenderedElement.ID())

		if err != nil {
			logger.Error("Error removing rendered element", zap.Error(err))
		}
	}

	for _, action := range actionsOther {
		switch a := action.(type) {
		case *renderActionCreate:
			logger.Debug("ExecuteRenderActions: creating new element", zap.Any("a", a))

			rendereredMessage, err := create(ctx, renderer, a, logger)

			if err != nil {
				return nil, err
			}

			result = append(result, rendereredMessage)
		case *renderActionKeep:
			logger.Debug("ExecuteRenderActions: keeping rendered element", zap.Any("a", a))

			if a.RenderedElement.renderedKind() == kindRenderedBotMessage && a.NewElement.OutcomingKind() == kindOutcomingTextMessage {
				rendereredMessage := &renderedBotMessage{
					OutcomingTextMessage: a.NewElement.(*outcomingTextMessage),
					Message:              a.RenderedElement.(*renderedBotMessage).Message,
				}

				result = append(result, rendereredMessage)
			}

		case *renderActionReplace:
			if a.RenderedElement.renderedKind() == kindRenderedBotMessage && a.NewElement.OutcomingKind() == kindOutcomingTextMessage {

				logger.Debug("ExecuteRenderActions: replacing rendered element", zap.Any("a", a.RenderedElement))

				outcoming := a.NewElement.(*outcomingTextMessage)
				renderedElement := a.RenderedElement.(*renderedBotMessage)

				logger.Debug("ExecuteRenderActions: replacing rendered element",
					zap.Any("outcoming", outcoming), zap.Any("renderedElement", renderedElement),
				)

				message, err := renderer.Message(ctx, &ChatRendererMessageProps{
					Text:          getOrText(outcoming.Text, emptyString),
					ReplyMarkup:   outcoming.ReplyMarkup(),
					TargetMessage: renderedElement.Message,
					RemoveTarget:  false,
				})

				if err != nil {
					return nil, err
				}

				rendereredMessage := &renderedBotMessage{
					OutcomingTextMessage: outcoming,
					Message:              message,
				}

				result = append(result, rendereredMessage)
			}
		}
	}

	for _, a := range actionsRemoveBot {
		logger.Debug("ExecuteRenderActions: removing rendered element", zap.Any("a", a))

		err := renderer.Delete(a.RenderedElement.ID())

		if err != nil {
			logger.Error("Error removing rendered element", zap.Error(err))
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
