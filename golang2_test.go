package tgbot_test

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflectState(t *testing.T) {

}

type S struct {
	s string
	a int
}

type Comp interface {
	Print()
}

func (s S) Print() {
	fmt.Println(s)
}

func printProps(c Comp) {
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

	var c Comp = s

	printProps(c)

}
