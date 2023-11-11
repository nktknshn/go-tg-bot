package tgbot

import "reflect"

func reflectCompProps[A any](comp Comp[A]) reflect.Value {
	t := reflect.TypeOf(comp)

	fs := make([]reflect.StructField, 0)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.Name == "state" {
			continue
		}

		fs = append(fs, f)
	}

	props := reflect.New(reflect.StructOf(fs)).Elem()
	props.Set(reflect.ValueOf(comp))

	return props
}

func reflectCompId[A any](comp Comp[A]) string {
	t := reflect.TypeOf(comp)
	return t.Name()
}
