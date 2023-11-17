package tgbot

import (
	"go.uber.org/zap"
)

func ComponentToElements[A any](comp Comp[A], logger *zap.Logger) []Element {

	logger.Debug("ComponentToElements", zap.Any("comp", comp))

	elements := make([]Element, 0)

	o := NewOutput[A]()

	comp.Render(o)

	for _, e := range o.Result {
		switch e := e.(type) {
		case *ElementComponent[A]:
			logger.Debug("Going into ElementComponent", zap.Any("comp", e.comp))
			compElements := ComponentToElements(e.comp, logger)
			elements = append(elements, compElements...)
		default:
			elements = append(elements, e)
		}
	}

	logger.Debug("ComponentToElements", Elements(elements).ZapField("elements"))

	return elements
}
