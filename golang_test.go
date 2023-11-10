package tgbot_test

import (
	"fmt"
	"testing"
)

func TestStructSlice(t *testing.T) {

	type S struct {
		ss []string
	}

	s := S{}
	fmt.Println(s)

	s.ss = append(s.ss, "a")

	fmt.Println(s)

}

type A struct {
	a string
}

func (a A) String() string {
	return fmt.Sprintf("A{a=%s}", a.a)
}

type B struct {
	b string
}

func (b *B) String() string {
	return fmt.Sprintf("B{b=%s}", b.b)
}

func TestStringMethod(t *testing.T) {
	// conclusion: method for pointer is only called when the variable is a pointer
	// method for value is called when the variable is a value or a pointer

	a := A{a: "a"}
	ap := &A{a: "ap"}
	b := B{b: "b"}
	bp := &B{b: "bp"}

	fmt.Println(a)
	fmt.Println(ap)
	fmt.Println(b)
	fmt.Println(bp)

}

func TestStructNested(t *testing.T) {

	type S struct {
		s string
	}

	type T struct {
		S
	}

	tt := T{
		S{
			s: "test",
		},
	}

	fmt.Println(tt.S.s)
	fmt.Println(tt.s)
}

func TestClosureCopy(t *testing.T) {

	type V struct{ v int }

	values := []V{{1}, {2}, {3}}
	callbacks := make([]func() (V, int), 0)

	for idx, v := range values {
		idx := idx
		v := v
		callbacks = append(callbacks, func() (V, int) { return v, idx })
	}

	for _, cb := range callbacks {
		fmt.Println(cb())
	}

}
