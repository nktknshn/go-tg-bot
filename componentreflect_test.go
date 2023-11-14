package tgbot_test

import (
	"reflect"
	"testing"

	tgbot "github.com/nktknshn/go-tg-bot"
)

type TestReflectCompCtxReqsApp1 struct {
	Flag1    bool
	CtxFlag2 bool
	CtxName  string
}

func (c *TestReflectCompCtxReqsApp1) Render(o tgbot.OO) {

}

func TestReflectCompCtxReqs(t *testing.T) {

	type Reqs struct {
		Flag2 bool
		Name  string
	}

	ctxReqs := tgbot.ReflectCompCtxReqs(&TestReflectCompCtxReqsApp1{})

	// fmt.Println(ctxReqs.FieldByName("Flag2"))

	// fmt.Println()
	if !ctxReqs.Type().ConvertibleTo(reflect.TypeOf(Reqs{})) {
		t.Fatal("Expected the Reqs")
	}

	// fmt.Println(Reqs(ctxReqs.Interface()))
	// fmt.Println(ctxReqs.Interface().(Reqs).Flag2)

}

type TestReflectTagApp1 struct {
	Flag1 bool
	Flag2 bool   `tgbot:"ctx"`
	Name  string `tgbot:"ctx"`
}

func (c *TestReflectTagApp1) Render(o tgbot.OO) {
}

func TestReflectTag1(t *testing.T) {
	reqs := tgbot.ReflectCompCtxReqsTags(&TestReflectTagApp1{})

	if reqs.Get("Flag2") != reflect.TypeOf(true) {
		t.Fatal("Expected bool")
	}

	if reqs.Get("Name") != reflect.TypeOf("") {
		t.Fatal("Expected string")
	}
}

type TestReflectTagApp2 struct {
}

func (c *TestReflectTagApp2) Render(o tgbot.OO) {
}

func TestReflectTag2(t *testing.T) {
	reqs := tgbot.ReflectCompCtxReqsTags(&TestReflectTagApp2{})

	if !reqs.IsEmpty() {
		t.Fatal("Expected empty")
	}
}
