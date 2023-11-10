package tgbot

import "go.uber.org/zap"

func ComponentToElements[A any](comp Comp[A]) []Element {

	elements := make([]Element, 0)

	o := newOutput[A]()

	comp.Render(o)

	for _, e := range o.result {
		switch e := e.(type) {
		case *ElementComponent[A]:
			compElements := ComponentToElements(e.comp)
			elements = append(elements, compElements...)
		default:
			elements = append(elements, e)
		}
	}

	globalLogger.Debug("ComponentToElements",
		zap.Any("elements", elements),
	)

	return elements
}
