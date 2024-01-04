package render

import (
	"slices"

	"github.com/nktknshn/go-tg-bot/tgbot/outcoming"
	"go.uber.org/zap"
)

func CreateRenderActions(renderedElements []RenderedElement, nextElements []outcoming.OutcomingMessage, logger *zap.Logger) []RenderAction {

	logger.Debug("CreateRenderActions",
		zap.Any("renderedElements", len(renderedElements)),
		zap.Any("nextElements", len(nextElements)),
	)

	actions := make([]RenderAction, 0)

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
		var n outcoming.OutcomingMessage

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

		if n == nil {
			// we are out of new elements to render so we can delete all remaining rendered elements
			result = slices.Delete(result, idx, idx+1)
			idx -= 1
			actions = append(actions, &renderActionRemove{RenderedElement: r})
		} else if r == nil {
			// we are out of rendered elements so we can create all remaining new elements
			actions = append(actions, &renderActionCreate{NewElement: n})
			continue
		} else if isRenderedEqualOutcoming(r, n) {
			actions = append(actions, &renderActionKeep{RenderedElement: r, NewElement: n})
		} else if slices.IndexFunc(renderedElements, func(re RenderedElement) bool { return isRenderedEqualOutcoming(re, n) }) > idx {
			// if we have the next outcoming element rendered somewhere else ahead of current rendered element
			// we can delete current rendered element
			result = slices.Delete(result, idx, idx+1)
			nextElements = append([]outcoming.OutcomingMessage{n}, nextElements...)
			idx -= 1
			actions = append(actions, &renderActionRemove{RenderedElement: r})
		} else {
			if r.CanReplace(n) {
				actions = append(actions, &renderActionReplace{RenderedElement: r, NewElement: n})
			} else {
				result = slices.Delete(result, idx, idx+1)
				nextElements = append([]outcoming.OutcomingMessage{n}, nextElements...)
				idx -= 1
				actions = append(actions, &renderActionRemove{RenderedElement: r})
			}
		}

		idx += 1

	}

	return actions
}
