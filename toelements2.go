package tgbot

import "go.uber.org/zap"

func createElements[A any](comp Comp[A], logger *zap.Logger) {
	o := newOutput[A]()

	comp.Render(o)

}
