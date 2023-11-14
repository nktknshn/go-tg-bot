package tgbot_test

import (
	"reflect"
	"testing"

	tgbot "github.com/nktknshn/go-tg-bot"
)

func TestCreateElementsContext(t *testing.T) {
	ctx := tgbot.NewCreateElementsContext()

	ctx.Add("A", 1)
	ctx.Add("B", "2")
	ctx.Add("C", true)

	q := tgbot.NewContextQueryFromMap(map[string]reflect.Type{
		"A": reflect.TypeOf(1),
		"B": reflect.TypeOf(""),
		"C": reflect.TypeOf(true),
	})

	res, err := ctx.Query(q)

	if err != nil {
		t.Fatal(err)
	}

	// fmt.Println(res.)
	t.Logf("%v", res)
}
