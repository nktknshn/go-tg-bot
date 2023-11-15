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

		fs = append(fs, f)
	}

	props := reflect.New(reflect.StructOf(fs)).Elem()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.Name == "State" {
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

	if ok != true {
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
