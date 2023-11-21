package tgbot

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
)

const localStateClosureName = "LocalStateClosure"

func reflectCompProps(comp Comp) any {

	t := reflect.TypeOf(comp)
	// .Elem()
	v := reflect.ValueOf(comp)
	// .Elem()

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	fs := make([]reflect.StructField, 0)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.Name == "State" {
			continue
		}

		if f.Name == "Context" {
			continue
		}

		fs = append(fs, f)
	}

	props := reflect.New(reflect.StructOf(fs)).Elem()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.Name == "State" {
			continue
		}

		if f.Name == "Context" {
			continue
		}

		// fmt.Println("Setting", f.Name)

		props.FieldByName(f.Name).Set(
			v.FieldByName(f.Name),
		)

	}
	// props.Set(reflect.ValueOf(comp).Elem())

	return props.Interface()
}

func reflectCompId(comp Comp) string {
	t := reflect.TypeOf(comp)
	return fmt.Sprintf("%v", t)
}

// sets the local state of the component if it has one defined
// returns a copy of the component where State.LocalStateClosure initialized to zero value
// and filled with values from state if any
func reflectCompLocalState(
	logger *zap.Logger,
	comp Comp,
	state CompState[any],
) Comp {
	logger.Debug("ReflectCompLocalState")

	wasapointer := false

	t := reflect.TypeOf(comp)
	v := reflect.ValueOf(comp)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
		wasapointer = true
	}

	compStateField, ok := t.FieldByName("State")

	if !ok {
		logger.Debug("Component doesn't use local state")
		return comp
	}

	compStateClosurePtrField, ok := compStateField.Type.FieldByName(localStateClosureName)

	if !ok {
		panic(fmt.Errorf("ReflectCompLocalState: %w", fmt.Errorf("component has State but doesn't have LocalStateClosure")))
	}

	compStateClosureType := compStateClosurePtrField.Type.Elem()

	stateValue := reflect.ValueOf(state)
	stateClosurePtr := stateValue.FieldByName(localStateClosureName)
	stateClosureValue := stateClosurePtr.Elem().FieldByName("Value")

	compCopy := reflect.New(t).Elem()

	for i := 0; i < compCopy.NumField(); i++ {
		compCopy.Field(i).Set(v.Field(i))
	}

	// copy index
	compCopy.FieldByName("State").FieldByName("Index").Set(
		stateValue.FieldByName("Index"),
	)

	if state.LocalStateClosure == nil {
		panic("state.LocalState is nil. Initialize zero closure before running the component")
	}

	newClosurePtr := reflect.New(compStateClosureType)

	compCopy.FieldByName("State").FieldByName(localStateClosureName).Set(
		newClosurePtr,
	)

	if !state.LocalStateClosure.Initialized {
		logger.Debug("No input state (component first run). Default will be used.")

		// copy ptr to the closure
		// component Render func fill use it to write the initial values

		return compCopy.Addr().Interface().(Comp)
	}

	logger.Debug("Filling component local state from input state",
		zap.Any("state", stateClosureValue.Interface()),
	)

	newClosurePtr.Elem().FieldByName("Initialized").SetBool(
		true,
	)

	newClosurePtr.Elem().FieldByName("Value").Set(
		reflect.ValueOf(stateClosureValue.Interface()),
	)

	if wasapointer {
		return compCopy.Addr().Interface().(Comp)
	}

	return compCopy.Interface().(Comp)
}

func (ucv usedContextValue) Interface() []any {
	result := make([]any, 0)

	for _, v := range ucv {
		result = append(result, v.Interface())
	}

	return result
}

func (ucv usedContextValue) Equal(other usedContextValue) bool {
	return reflect.DeepEqual(ucv.Interface(), other.Interface())
}

func reflectDerefValue(v reflect.Value) reflect.Value {

	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}

	return v
}

func reflectHasState(comp Comp) bool {
	t := reflect.TypeOf(comp)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	_, ok := t.FieldByName("State")

	return ok
}

// Try to find `Context` field.
// If found fill it with the global context values returning a new component.
// Returns the new component and a pointer to the used context value (if any).
// If the component has a `Selector` method, it will be called to get the context value.
func reflectTypedContext[C any](comp Comp, globalContext C) (Comp, *usedContextValue) {
	var wasapointer = false

	t := reflect.TypeOf(comp)
	v := reflect.ValueOf(comp)

	// dereference pointer
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
		wasapointer = true
	}

	_, hasContext := t.FieldByName("Context")

	if !hasContext {
		return comp, nil
	}

	compCopy := reflect.New(t).Elem()

	for i := 0; i < t.NumField(); i++ {
		compCopy.Field(i).Set(v.Field(i))
	}

	compCopyCtx := compCopy.FieldByName("Context")

	compCopyCtx.Set(reflect.ValueOf(globalContext))

	var usedContextValue *usedContextValue = reflectTypedContextSelect[C](comp, globalContext)

	if wasapointer {
		return compCopy.Addr().Interface().(Comp), usedContextValue
	}

	return compCopy.Interface().(Comp), usedContextValue
}

// Returns a pointer to the used part of the global context (if any).
func reflectTypedContextSelect[C any](comp Comp, globalContext C) *usedContextValue {

	t := reflect.TypeOf(comp)

	// dereference pointer
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	_, hasContext := t.FieldByName("Context")
	selector, hasSelector := t.MethodByName("Selector")

	if !hasContext {
		return nil
	}

	var usedContextValue usedContextValue

	compCopy := reflect.New(t).Elem()

	compCopyCtx := compCopy.FieldByName("Context")

	ctxValue := globalContext

	compCopyCtx.Set(reflect.ValueOf(ctxValue))

	if hasSelector {
		usedContextValue = selector.Func.Call([]reflect.Value{compCopy})
	} else {
		usedContextValue = []reflect.Value{reflect.ValueOf(globalContext)}
	}

	return &usedContextValue
}
