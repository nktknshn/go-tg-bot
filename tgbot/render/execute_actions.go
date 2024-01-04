package render

import (
	"context"

	"github.com/nktknshn/go-tg-bot/helpers"
	"github.com/nktknshn/go-tg-bot/tgbot/outcoming"
	"go.uber.org/zap"
)

type ExecuteRenderActionsProps struct {
	Logger      *zap.Logger
	LoggerTrace *zap.Logger
}

func (props ExecuteRenderActionsProps) Loggers() ExecuteRenderActionsProps {
	if props.Logger == nil {
		props.Logger = zap.NewNop()
	}

	if props.LoggerTrace == nil {
		props.LoggerTrace = zap.NewNop()
	}

	return ExecuteRenderActionsProps{
		Logger:      props.Logger,
		LoggerTrace: props.LoggerTrace,
	}
}

// Takes actions and applies them to the renderer
func ExecuteRenderActions(ctx context.Context, renderer ChatRenderer, actions []RenderAction, props ExecuteRenderActionsProps) ([]RenderedElement, error) {

	loggers := props.Loggers()

	result := make([]RenderedElement, 0)
	actionsRemove := make([]renderActionRemove, 0)
	actionsRemoveBot := make([]renderActionRemove, 0)
	actionsRemoveUser := make([]renderActionRemove, 0)

	actionsOther := make([]RenderAction, 0)

	for _, action := range actions {
		switch a := action.(type) {
		case *renderActionRemove:
			if _, ok := a.RenderedElement.(*RenderedBotMessage); ok {
				actionsRemoveBot = append(actionsRemoveBot, *a)
			} else if _, ok := a.RenderedElement.(*RenderedUserMessage); ok {
				actionsRemoveUser = append(actionsRemoveUser, *a)
			}
			actionsRemove = append(actionsRemove, *a)
		default:
			actionsOther = append(actionsOther, a)
		}
	}

	for _, a := range actionsRemoveUser {
		loggers.Logger.Debug("ExecuteRenderActions: removing rendered element", zap.Any("a", a))

		err := renderer.Delete(a.RenderedElement.ID())

		if err != nil {
			loggers.Logger.Error("Error removing rendered element", zap.Error(err))
		}
	}

	for _, action := range actionsOther {
		switch a := action.(type) {
		case *renderActionCreate:
			loggers.Logger.Debug("ExecuteRenderActions: creating new element", zap.Any("a", a))

			rendereredMessage, err := create(ctx, renderer, *a, loggers)

			if err != nil {
				return nil, err
			}

			result = append(result, rendereredMessage)
		case *renderActionKeep:
			loggers.Logger.Debug("ExecuteRenderActions: keeping rendered element", zap.Any("a", a))

			if r, ok := a.RenderedElement.(*RenderedBotMessage); ok {
				if n, ok := a.NewElement.(*outcoming.OutcomingTextMessage); ok {
					rendereredMessage := &RenderedBotMessage{
						Message:              r.Message,
						OutcomingTextMessage: n,
					}

					result = append(result, rendereredMessage)
				}
			}

			// if a.RenderedElement.RenderedKind() == KindRenderedBotMessage && a.NewElement.OutcomingKind() == outcoming.KindOutcomingTextMessage {
			// 	rendereredMessage := &RenderedBotMessage{
			// 		OutcomingTextMessage: a.NewElement.(*outcoming.OutcomingTextMessage),
			// 		Message:              a.RenderedElement.(*RenderedBotMessage).Message,
			// 	}

			// 	result = append(result, rendereredMessage)
			// }

		case *renderActionReplace:
			// if a.RenderedElement.RenderedKind() == KindRenderedBotMessage && a.NewElement.OutcomingKind() == outcoming.KindOutcomingTextMessage {
			if r, ok := a.RenderedElement.(*RenderedBotMessage); ok {
				if n, ok := a.NewElement.(*outcoming.OutcomingTextMessage); ok {

					outcoming := n
					renderedElement := r

					loggers.Logger.Debug("ExecuteRenderActions: replacing rendered element",
						zap.Any("outcoming", outcoming), zap.Any("renderedElement", renderedElement),
					)

					message, err := renderer.Message(ctx, &ChatRendererMessageProps{
						Text:          helpers.GetOrText(outcoming.Text, emptyString),
						ReplyMarkup:   outcoming.ReplyMarkup(),
						TargetMessage: renderedElement.Message,
						RemoveTarget:  false,
					})

					if err != nil {
						return nil, err
					}

					rendereredMessage := &RenderedBotMessage{
						OutcomingTextMessage: outcoming,
						Message:              message,
					}

					result = append(result, rendereredMessage)
				}
			}
		}
	}

	for _, a := range actionsRemoveBot {
		loggers.Logger.Debug("ExecuteRenderActions: removing rendered element", zap.Any("a", a))

		err := renderer.Delete(a.RenderedElement.ID())

		if err != nil {
			loggers.Logger.Error("Error removing rendered element", zap.Error(err))
		}
	}

	return result, nil
}

const emptyString = "<empty>"

func create(ctx context.Context, renderer ChatRenderer, action renderActionCreate, props ExecuteRenderActionsProps) (RenderedElement, error) {

	switch a := action.NewElement.(type) {
	case *outcoming.OutcomingTextMessage:

		message, err := renderer.Message(ctx, &ChatRendererMessageProps{
			Text:        helpers.GetOrText(a.Text, emptyString),
			ReplyMarkup: a.ReplyMarkup(),
		})

		if err != nil {
			return nil, err
		}

		return &RenderedBotMessage{
			Message:              message,
			OutcomingTextMessage: a,
		}, nil

	case *outcoming.OutcomingUserMessage:
		return &RenderedUserMessage{
			MessageID:            a.ElementUserMessage.MessageID,
			OutcomingUserMessage: a,
		}, nil
	// TODO
	// case *OutcomingFileMessage:
	// 	return renderer.File(a.ElementFile)
	// case *OutcomingPhotoGroupMessage:
	// 	return renderer.PhotoGroup(a.ElementPhotoGroup)
	default:
		props.Logger.Error("create: unsupported outcoming message type", zap.Any("a", a))
	}

	return nil, nil

}
