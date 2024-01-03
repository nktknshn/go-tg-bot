package tgbot

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotd/td/tg"
)

func Check(
	t *testing.T,
	renderedElements []RenderedElement,
	outcomingMessages []outcomingMessage,
	expected []renderActionType,
) {
	res0 := getRenderActions(renderedElements, outcomingMessages)

	if len(res0) != len(expected) {
		for i, v := range res0 {
			fmt.Printf("res0[%d] = %v\n", i, v.RenderActionKind())
		}

		t.Fatalf("len(res0) != len(expected). %d != %d", len(res0), len(expected))
	}

	for i, v := range res0 {
		e := expected[i]

		if v.RenderActionKind() != e.RenderActionKind() {
			t.Logf("expected  %v", e.RenderActionKind())
			t.Logf("received %v", v.RenderActionKind())
			t.Fatalf("res0[%d].RenderActionKind() != expected[%d].RenderActionKind()", i, i)
		}

		if v.RenderActionKind() == kindRenderActionCreate {

			check := v.(*renderActionCreate).NewElement.Equal(e.(*renderActionCreate).NewElement)

			if !check {
				t.Fatalf("res0[%d].NewElement != expected[%d].NewElement", i, i)
			}
		}

		if v.RenderActionKind() == kindRenderActionKeep {

			check := v.(*renderActionKeep).NewElement.Equal(e.(*renderActionKeep).NewElement)
			check = check && v.(*renderActionKeep).RenderedElement.Equal(e.(*renderActionKeep).RenderedElement)

			if !check {
				t.Fatalf("res0[%d].NewElement != expected[%d].NewElement", i, i)
			}
		}

		if v.RenderActionKind() == kindRenderActionReplace {

			check := v.(*renderActionReplace).NewElement.Equal(e.(*renderActionReplace).NewElement)
			check = check && v.(*renderActionReplace).RenderedElement.Equal(e.(*renderActionReplace).RenderedElement)

			if !check {
				t.Fatalf("res0[%d].NewElement != expected[%d].NewElement", i, i)
			}
		}

		if v.RenderActionKind() == kindRenderActionRemove {

			check := v.(*renderActionRemove).RenderedElement.Equal(e.(*renderActionRemove).RenderedElement)

			if !check {
				t.Fatalf("res0[%d].RenderedElement != expected[%d].RenderedElement", i, i)
			}
		}

		t.Logf("%v = %v", v.RenderActionKind(), e.RenderActionKind())
	}
}

var (
	m1 = newOutcomingTextMessage("message 1")
	m2 = newOutcomingTextMessage("message 2")
	m3 = newOutcomingTextMessage("message 3")
	m4 = newOutcomingTextMessage("message 4")

	rm1 = &renderedBotMessage{
		OutcomingTextMessage: m1,
		Message:              &tg.Message{},
	}
	rm2 = &renderedBotMessage{
		OutcomingTextMessage: m2,
		Message:              &tg.Message{},
	}
	rm3 = &renderedBotMessage{
		OutcomingTextMessage: m3,
		Message:              &tg.Message{},
	}
	rm4 = &renderedBotMessage{
		OutcomingTextMessage: m4,
		Message:              &tg.Message{},
	}
)

func TestGetRenderActionsInsertedMiddle(t *testing.T) {
	Check(t,
		[]RenderedElement{rm1, rm2, rm3},
		[]outcomingMessage{m1, m2, m4, m3},
		[]renderActionType{
			&renderActionKeep{RenderedElement: rm1, NewElement: m1},
			&renderActionKeep{RenderedElement: rm2, NewElement: m2},
			&renderActionReplace{RenderedElement: rm3, NewElement: m4},
			&renderActionCreate{NewElement: m3},
		})

}

func TestGetRenderActionsInsertedFirst(t *testing.T) {
	Check(t,
		[]RenderedElement{rm1, rm2, rm3},
		[]outcomingMessage{m4, m1, m2, m3},
		[]renderActionType{
			&renderActionReplace{RenderedElement: rm1, NewElement: m4},
			&renderActionReplace{RenderedElement: rm2, NewElement: m1},
			&renderActionReplace{RenderedElement: rm3, NewElement: m2},
			&renderActionCreate{NewElement: m3},
		})

}

func TestGetRenderActionsBasic(t *testing.T) {
	Check(t,
		[]RenderedElement{},
		[]outcomingMessage{m1},
		[]renderActionType{
			&renderActionCreate{NewElement: m1},
		})

	Check(t,
		[]RenderedElement{},
		[]outcomingMessage{m1, m2, m3},
		[]renderActionType{
			&renderActionCreate{NewElement: m1},
			&renderActionCreate{NewElement: m2},
			&renderActionCreate{NewElement: m3},
		})

	Check(t,
		[]RenderedElement{
			&renderedBotMessage{
				OutcomingTextMessage: m2,
			},
		},
		[]outcomingMessage{m2},
		[]renderActionType{
			&renderActionKeep{
				RenderedElement: &renderedBotMessage{
					OutcomingTextMessage: m2,
				},
				NewElement: m2,
			},
		})

	// rm1 is supposed to be replaced with m2
	Check(t,
		[]RenderedElement{rm1},
		[]outcomingMessage{m2},
		[]renderActionType{
			&renderActionReplace{
				RenderedElement: rm1,
				NewElement:      m2,
			},
		})

	// rm1 is supposed to be removed
	Check(t,
		[]RenderedElement{rm1},
		[]outcomingMessage{},
		[]renderActionType{
			&renderActionRemove{
				RenderedElement: rm1,
			},
		})

}

type MockRenderer struct {
	OutcomingMessages []outcomingMessage
}

func (mr *MockRenderer) Message(ctx context.Context, props *ChatRendererMessageProps) (*tg.Message, error) {
	return &tg.Message{}, nil
}

func (mr *MockRenderer) Delete(messageId int) error {
	return nil
}

func TestCreate(t *testing.T) {
	re := &MockRenderer{}

	m := newOutcomingTextMessage("message 1")
	b1 := newButton("button 1", func() any { return 1 }, "button 1", false, false)
	m.AddButton(b1)

	ExecuteRenderActions(
		context.Background(),
		re,
		[]renderActionType{
			&renderActionCreate{
				NewElement: m,
			},
		},
		DevLogger(),
	)
}
