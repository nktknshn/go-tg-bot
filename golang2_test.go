package tgbot_test

import (
	"fmt"
	"reflect"
	"testing"
)

type S struct {
	s string
	a int
}

type TestReflect1C interface {
	Print()
}

func (s S) Print() {
	fmt.Println(s)
}

func printProps(c TestReflect1C) {
	name := reflect.TypeOf(c).Name()
	fmt.Println(name)

	pkg := reflect.TypeOf(c).PkgPath()
	fmt.Println(pkg)

	t := reflect.TypeOf(c)

	// copy c into c2
	c2 := reflect.New(t).Elem()
	c2.Set(reflect.ValueOf(c))

	fmt.Println(c2)
}

func TestReflect1(t *testing.T) {

	s := S{s: "s", a: 1}

	var c TestReflect1C = s

	printProps(c)

}
