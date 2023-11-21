package tgbot

import (
	"fmt"
	"testing"
)

type TestReflectTypedContext1Context struct {
	value1 int
	value2 bool
}

type TestReflectTypedContext1Comp struct {
	Context TestReflectTypedContext1Context
}

func (c TestReflectTypedContext1Comp) Selector() (int, bool) {
	return c.Context.value1, c.Context.value2
}

func (c TestReflectTypedContext1Comp) Render(o O) {
	v1, v2 := c.Selector()

	o.Messagef("%v %v", v1, v2)
}

func TestReflectTypedContext1(t *testing.T) {
	ctx := TestReflectTypedContext1Context{
		value1: 1,
		value2: true,
	}

	comp := (Comp)(&TestReflectTypedContext1Comp{})

	comp, userContext := reflectTypedContext(comp, ctx)

	fmt.Println("Used context:", (*userContext)[1].Interface())

	o := newOutput()
	comp.Render(o)

	if o.Result[0].(*elementMessage).Text != "1 true" {
		t.Fatal("Expected 1 true")
	}
}
