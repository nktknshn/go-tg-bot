package tgbot

import "reflect"

func reflectCompProps[A any](comp Comp[A]) reflect.Value {
	t := reflect.TypeOf(comp).Elem()

	fs := make([]reflect.StructField, 0)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// if f.Name == "State" {
		// 	continue
		// }

		fs = append(fs, f)
	}

	props := reflect.New(reflect.StructOf(fs)).Elem()
	props.Set(reflect.ValueOf(comp).Elem())

	return props
}

func reflectCompId[A any](comp Comp[A]) string {
	t := reflect.TypeOf(comp)
	return t.Name()
}
