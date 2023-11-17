package tgbot

import (
	"fmt"
	"reflect"
	"strings"
)

func reflectCompProps[A any](comp Comp[A]) any {
	t := reflect.TypeOf(comp).Elem()
	v := reflect.ValueOf(comp).Elem()

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

// sets the local state of the component if it has one defined
func ReflectCompLocalState[A any](comp Comp[A], ls State[any]) Comp[A] {

	// fmt.Println("ReflectCompLocalState")
	// isPointer := false
	// TODO

	t := reflect.TypeOf(comp).Elem()

	// fmt.Println("t: ", t)

	_, ok := t.FieldByName("State")

	if !ok {
		// fmt.Println("Component doesn't use local state")
		return comp
	}

	v := reflect.ValueOf(comp).Elem()
	// vp := reflect.ValueOf(&comp).Elem()

	// fmt.Println("v: ", v)
	// fmt.Println("v.Type(): ", v.Type())
	// fmt.Println("vp: ", vp)
	// fmt.Println("vp.Type(): ", vp.Type())

	// fmt.Println("vp.Elem(): ", vp.Elem())
	// fmt.Println("vp.Elem().Type(): ", vp.Elem().Type())
	// fmt.Println("vp.Elem().Elem().Type(): ", vp.Elem().Elem().Type())

	// sf := vp.Elem().Elem().FieldByName(stateField.Name)

	vls := reflect.ValueOf(ls)

	// fmt.Printf("sf: %v\n", sf.Type())
	// fmt.Printf("vls: %v\n", vls.Type())

	// for i := 0; i < vls.NumField(); i++ {
	// 	fmt.Printf("vls.Field(%v): %v\n", i, vls.Field(i).Type().Name())
	// }

	vlsValue := vls.FieldByName("LocalState").FieldByName("Value")

	// fmt.Println("vlsValue: ", reflect.TypeOf(vlsValue))
	// fmt.Println("vlsValue: ", reflect.TypeOf(vlsValue.Interface()))

	// fmt.Println("sf.CanSet(): ", sf.CanSet())

	nt := reflect.New(t).Elem()

	// fmt.Println("nt.Type(): ", nt.Type())

	// ntf := nt.Interface()
	// nts := reflect.ValueOf(&ntf).Elem()

	// fmt.Println("nt.CanSet(): ", nt.CanSet())
	// fmt.Println("nt.CanSet(): ", nts.CanAddr())

	for i := 0; i < nt.NumField(); i++ {
		// copy props
		nt.Field(i).Set(v.Field(i))
		// fmt.Printf("nt.Field(%v): %v\n", i, nt.Field(i))
	}

	// fmt.Println("vls.Type()", vls.Type())
	// fmt.Println("nt.state", nt.FieldByName("State").Type())

	// fmt.Println(ls.Index)

	// set state index

	// stateType, _ := t.FieldByName("State")
	// newState := reflect.New(stateType.Type).Elem()

	nt.FieldByName("State").FieldByName("Index").Set(
		vls.FieldByName("Index"),
	)

	if !ls.LocalState.Initialized {
		// initialize local state
		// fmt.Println("No input state. Default will be used")
		// fmt.Println(nt.FieldByName("State"))

		return nt.Addr().Interface().(Comp[A])
	}

	nt.FieldByName("State").FieldByName("LocalState").FieldByName("Value").Set(
		reflect.ValueOf(vlsValue.Interface()),
	)
	nt.FieldByName("State").FieldByName("LocalState").FieldByName("Initialized").SetBool(
		true,
	)

	// fmt.Println("nt.value", nt.FieldByName("State").FieldByName("LocalState").FieldByName("Value"))
	// fmt.Println("nt.init", nt.FieldByName("State").FieldByName("LocalState").FieldByName("Initialized"))

	return nt.Addr().Interface().(Comp[A])
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

// Try to find `Context` field.
// If found fill it with the global context values returning a new component.
// Returns the new component and a pointer to the used context value (if any).
// If the component has a `Selector` method, it will be called to get the context value.
func ReflectTypedContext[A any, C any](comp Comp[A], globalContext C) (Comp[A], *UsedContextValue) {
	var wasapointer = false

	t := reflect.TypeOf(comp)

	// dereference pointer
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		wasapointer = true
	}

	_, hasContext := t.FieldByName("Context")

	if !hasContext {
		return comp, nil
	}

	compCopy := reflect.New(t).Elem()

	for i := 0; i < t.NumField(); i++ {
		compCopy.Field(i).Set(reflect.ValueOf(comp).Elem().Field(i))
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
