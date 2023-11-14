package tgbot

import (
	"fmt"
	"reflect"
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

// Returns a struct that will be used to reqeust the global context
func reflectCompCtxReqs[A any](comp Comp[A]) {

}
