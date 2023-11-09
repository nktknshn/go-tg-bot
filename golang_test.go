package tgbot_test

import (
	"fmt"
	"testing"
)

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
