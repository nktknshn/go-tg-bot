package tgbot

import (
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/zap"
)

func reflectCompProps[A any](comp Comp[A]) any {

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

func reflectCompId[A any](comp Comp[A]) string {
	t := reflect.TypeOf(comp)
	return fmt.Sprintf("%v", t)
}

const LocalStateClosureName = "LocalStateClosure"

// sets the local state of the component if it has one defined
// returns a copy of the component where State.LocalStateClosure initialized to zero value
// and filled with values from state if any
func ReflectCompLocalState[A any](
	logger *zap.Logger,
	comp Comp[A],
	state State[any],
) Comp[A] {
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

	compStateClosurePtrField, ok := compStateField.Type.FieldByName(LocalStateClosureName)

	if !ok {
		panic(fmt.Errorf("ReflectCompLocalState: %w", fmt.Errorf("Component has State but doesn't have LocalStateClosure")))
	}

	compStateClosureType := compStateClosurePtrField.Type.Elem()

	stateValue := reflect.ValueOf(state)
	stateClosurePtr := stateValue.FieldByName(LocalStateClosureName)
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

	compCopy.FieldByName("State").FieldByName(LocalStateClosureName).Set(
		newClosurePtr,
	)

	if !state.LocalStateClosure.Initialized {
		logger.Debug("No input state (component first run). Default will be used.")

		// copy ptr to the closure
		// component Render func fill use it to write the initial values

		return compCopy.Addr().Interface().(Comp[A])
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
		return compCopy.Addr().Interface().(Comp[A])
	}

	return compCopy.Interface().(Comp[A])
}

// Returns a struct that will be used to request the global context
// Field that starts with Ctx will be queried
// comp is a pointer
func ReflectCompCtxReqs[A any](comp Comp[A]) reflect.Value {
	var prefix = "Ctx"
	t := reflect.TypeOf(comp)
	v := reflect.ValueOf(comp)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		v = v.Elem()
	}

	fs := make([]reflect.StructField, 0)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// fmt.Println("field ", f.Name)

		if !strings.HasPrefix(f.Name, prefix) {
			// fmt.Println("no prefix ", prefix)
			continue
		}

		fs = append(fs, reflect.StructField{
			Name:      f.Name[len(prefix):],
			PkgPath:   f.PkgPath,
			Type:      f.Type,
			Tag:       f.Tag,
			Index:     f.Index,
			Offset:    f.Offset,
			Anonymous: f.Anonymous,
		})
	}

	reqStructType := reflect.StructOf(fs)
	reqStruct := reflect.New(reqStructType)

	for _, f := range fs {
		reqStruct.Elem().FieldByName(f.Name).Set(v.FieldByName(prefix + f.Name))
	}

	return reqStruct.Elem()
}

type ContextQuery struct {
	reflect.Type
}

func (r ContextQuery) IsEmpty() bool {
	return r.NumField() == 0
}

func (r ContextQuery) Get(key string) reflect.Type {
	t, ok := r.FieldByName(key)

	if !ok {
		panic(ok)
	}

	return t.Type
}

func (r ContextQuery) String() string {
	result := ""

	for i := 0; i < r.NumField(); i++ {
		result += fmt.Sprintf("%s: %v\n", r.Field(i).Name, r.Field(i))
	}

	return result
}

func ReflectCompCtxReqsTags[A any](comp Comp[A]) ContextQuery {
	var tag = "tgbot"
	var tagValue = "ctx"

	t := reflect.TypeOf(comp)
	// v := reflect.ValueOf(comp)

	// dereference pointer
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		// v = v.Elem()
	}

	fs := make([]reflect.StructField, 0)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// fmt.Println("field ", f.Name)

		if f.Tag.Get(tag) != tagValue {
			continue
		}

		fs = append(fs, f)

	}

	reqStructType := reflect.StructOf(fs)
	// reqStruct := reflect.New(reqStructType)

	// for _, f := range fs {
	// 	reqStruct.Elem().FieldByName(f.Name).Set(v.FieldByName(f.Name))
	// }

	return ContextQuery{reqStructType}
}

func ReflectContextQueryResultGet[A any](comp Comp[A], globalContext GlobalContext) *ContextQueryResult {
	q := ReflectCompCtxReqsTags[A](comp)

	if q.IsEmpty() {
		return nil
	}

	res, err := globalContext.Query(q)

	if err != nil {
		panic(fmt.Errorf("ReflectContextQueryResultGet: %w", err))
	}

	return res
}

func ReflectContextQueryResultSet[A any](comp Comp[A], cqr *ContextQueryResult) Comp[A] {

	var wasapointer = false

	if cqr == nil {
		return comp
	}

	t := reflect.TypeOf(comp)
	v := reflect.ValueOf(comp)

	// dereference pointer
	if t.Kind() == reflect.Pointer {
		wasapointer = true
		t = t.Elem()
		v = v.Elem()
	}

	q := ReflectCompCtxReqsTags[A](comp)

	if q.IsEmpty() {
		//
		return comp
	}

	// component copy
	nt := reflect.New(t).Elem()

	for i := 0; i < t.NumField(); i++ {
		nt.Field(i).Set(v.Field(i))
	}

	for i := 0; i < q.NumField(); i++ {
		f := q.Field(i)

		// fmt.Println("field ", f.Name)

		if f.Tag.Get("tgbot") != "ctx" {
			continue
		}

		nt.FieldByName(f.Name).Set(cqr.Get(f.Name))
	}

	if wasapointer {
		return nt.Addr().Interface().(Comp[A])
	}

	return nt.Interface().(Comp[A])
}

type UsedContextValue []reflect.Value

func (ucv UsedContextValue) Interface() []any {
	result := make([]any, 0)

	for _, v := range ucv {
		result = append(result, v.Interface())
	}

	return result
}

func (ucv UsedContextValue) Equal(other UsedContextValue) bool {
	return reflect.DeepEqual(ucv.Interface(), other.Interface())
}

func ReflectDerefValue(v reflect.Value) reflect.Value {

	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}

	return v
}

func ReflectHasState[A any](comp Comp[A]) bool {
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
func ReflectTypedContext[A any, C any](comp Comp[A], globalContext C) (Comp[A], *UsedContextValue) {
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

	var usedContextValue *UsedContextValue = ReflectTypedContextSelect[A, C](comp, globalContext)

	if wasapointer {
		return compCopy.Addr().Interface().(Comp[A]), usedContextValue
	}

	return compCopy.Interface().(Comp[A]), usedContextValue
}

// Returns a pointer to the used part of the global context (if any).
func ReflectTypedContextSelect[A any, C any](comp Comp[A], globalContext C) *UsedContextValue {

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

	var usedContextValue UsedContextValue

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

func ReflectStructName(any any) string {
	return reflect.TypeOf(any).String()
}
