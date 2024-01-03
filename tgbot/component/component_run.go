package component

import (
	"reflect"

	"go.uber.org/zap"
)

func runComponent(
	logger *zap.Logger,
	comp Comp,
	gc GlobalContext[any],
	state CompState[any],
) ([]AnyElement, localStateClosure[any], *usedContextValue) {

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

	if !reflectHasState(comp) {
		logger.Debug("Component doesn't have state")
		return o.Result, *state.LocalStateClosure, usedContextValue
	}

	ls := reflectDerefValue(reflect.ValueOf(comp)).
		FieldByName("State")

	ls = ls.
		FieldByName(localStateClosureName).
		Elem()

	vi := ls.FieldByName("Initialized")
	vv := ls.FieldByName("Value")

	return o.Result, localStateClosure[any]{
		Initialized: vi.Bool(),
		Value:       vv.Interface(),
	}, usedContextValue

}
