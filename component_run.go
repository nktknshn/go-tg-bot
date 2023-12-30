package tgbot

import (
	"reflect"

	"go.uber.org/zap"
)

func runComponent(
	logger *zap.Logger,
	comp Comp,
	gc globalContext[any],
	state CompState[any],
) ([]anyElement, localStateClosure[any], *usedContextValue) {

	logger.Debug("RunComponent",
		zap.String("compId", reflectCompId(comp)),
		// zap.String("globalContext", globalContext),
		zap.Bool("isPointer", reflect.TypeOf(comp).Kind() == reflect.Ptr),
	)

	comp = reflectCompLocalState(logger, comp, state)

	comp, usedContextValue := reflectTypedContext(comp, gc.Get())

	o := newOutput()
	comp.Render(o)

	logger.Debug("Component rendered", zap.Any("comp", comp))

	ls := reflectDerefValue(reflect.ValueOf(comp)).
		FieldByName("State").
		FieldByName(localStateClosureName).
		Elem()

	vi := ls.FieldByName("Initialized")
	vv := ls.FieldByName("Value")

	// fmt.Println("Initialized", vi)
	// fmt.Println("Value", vv)

	if !reflectHasState(comp) {
		logger.Debug("Component doesn't have state")
		return o.Result, *state.LocalStateClosure, usedContextValue
	}

	return o.Result, localStateClosure[any]{
		Initialized: vi.Bool(),
		Value:       vv.Interface(),
	}, usedContextValue

}
