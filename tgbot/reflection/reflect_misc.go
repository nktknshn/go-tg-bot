package reflection

import "reflect"

func ReflectStructName(any any) string {
	return reflect.TypeOf(any).String()
}
