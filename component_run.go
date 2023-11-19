package tgbot

import (
	"reflect"

	"go.uber.org/zap"
)

func RunComponent[A any](
	logger *zap.Logger,
	comp Comp[A],
	globalContext GlobalContextTyped[any],
	state State[any],
) ([]Element, LocalStateClosure[any], *UsedContextValue) {

	logger.Debug("RunComponent",
		zap.String("compId", reflectCompId[A](comp)),
		// zap.String("globalContext", globalContext),
		zap.Bool("isPointer", reflect.TypeOf(comp).Kind() == reflect.Ptr),
	)

	comp = ReflectCompLocalState[A](logger, comp, state)

	comp, usedContextValue := ReflectTypedContext[A](comp, globalContext.Get())

	o := NewOutput[A]()
	comp.Render(o)

	ls := ReflectDerefValue(reflect.ValueOf(comp)).FieldByName("State").FieldByName(LocalStateClosureName).Elem()

	vi := ls.FieldByName("Initialized")
	vv := ls.FieldByName("Value")

	// fmt.Println("Initialized", vi)
	// fmt.Println("Value", vv)

	if !ReflectHasState(comp) {
		logger.Debug("Component doesn't have state")
		return o.Result, *state.LocalStateClosure, usedContextValue
	}

	return o.Result, LocalStateClosure[any]{
		Initialized: vi.Bool(),
		Value:       vv.Interface(),
	}, usedContextValue

}
