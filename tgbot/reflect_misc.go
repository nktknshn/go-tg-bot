package tgbot

import "reflect"

func reflectStructName(any any) string {
	return reflect.TypeOf(any).String()
}
