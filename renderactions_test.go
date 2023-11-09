package tgbot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-telegram/bot/models"
	tgbot "github.com/nktknshn/go-tg-bot"
)

func Check(
	t *testing.T,
	renderedElements []tgbot.RenderedElement,
	outcomingMessages []tgbot.OutcomingMessage,
	expected []tgbot.RenderActionType,
) {
	res0 := tgbot.GetRenderActions[any](renderedElements, outcomingMessages)

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

		if v.RenderActionKind() == tgbot.KindRenderActionCreate {

			check := v.(*tgbot.RenderActionCreate).NewElement.Equal(e.(*tgbot.RenderActionCreate).NewElement)

			if !check {
				t.Fatalf("res0[%d].NewElement != expected[%d].NewElement", i, i)
			}
		}

		if v.RenderActionKind() == tgbot.KindRenderActionKeep {

			check := v.(*tgbot.RenderActionKeep).NewElement.Equal(e.(*tgbot.RenderActionKeep).NewElement)
			check = check && v.(*tgbot.RenderActionKeep).RenderedElement.Equal(e.(*tgbot.RenderActionKeep).RenderedElement)

			if !check {
				t.Fatalf("res0[%d].NewElement != expected[%d].NewElement", i, i)
			}
		}

		if v.RenderActionKind() == tgbot.KindRenderActionReplace {

			check := v.(*tgbot.RenderActionReplace).NewElement.Equal(e.(*tgbot.RenderActionReplace).NewElement)
			check = check && v.(*tgbot.RenderActionReplace).RenderedElement.Equal(e.(*tgbot.RenderActionReplace).RenderedElement)

			if !check {
				t.Fatalf("res0[%d].NewElement != expected[%d].NewElement", i, i)
			}
		}

		if v.RenderActionKind() == tgbot.KindRenderActionRemove {

			check := v.(*tgbot.RenderActionRemove).RenderedElement.Equal(e.(*tgbot.RenderActionRemove).RenderedElement)

			if !check {
				t.Fatalf("res0[%d].RenderedElement != expected[%d].RenderedElement", i, i)
			}
		}

		t.Logf("%v = %v", v.RenderActionKind(), e.RenderActionKind())
	}
}

var (
	m1 = tgbot.NewOutcomingTextMessage[any]("message 1")
	m2 = tgbot.NewOutcomingTextMessage[any]("message 2")
	m3 = tgbot.NewOutcomingTextMessage[any]("message 3")
	m4 = tgbot.NewOutcomingTextMessage[any]("message 4")

	rm1 = &tgbot.RenderedBotMessage[any]{
		OutcomingTextMessage: m1,
		Message:              &models.Message{},
	}
	rm2 = &tgbot.RenderedBotMessage[any]{
		OutcomingTextMessage: m2,
		Message:              &models.Message{},
	}
	rm3 = &tgbot.RenderedBotMessage[any]{
		OutcomingTextMessage: m3,
		Message:              &models.Message{},
	}
	rm4 = &tgbot.RenderedBotMessage[any]{
		OutcomingTextMessage: m4,
		Message:              &models.Message{},
	}
)

func TestGetRenderActionsInsertedMiddle(t *testing.T) {
	Check(t,
		[]tgbot.RenderedElement{rm1, rm2, rm3},
		[]tgbot.OutcomingMessage{m1, m2, m4, m3},
		[]tgbot.RenderActionType{
			&tgbot.RenderActionKeep{RenderedElement: rm1, NewElement: m1},
			&tgbot.RenderActionKeep{RenderedElement: rm2, NewElement: m2},
			&tgbot.RenderActionReplace{RenderedElement: rm3, NewElement: m4},
			&tgbot.RenderActionCreate{NewElement: m3},
		})

}

func TestGetRenderActionsInsertedFirst(t *testing.T) {
	Check(t,
		[]tgbot.RenderedElement{rm1, rm2, rm3},
		[]tgbot.OutcomingMessage{m4, m1, m2, m3},
		[]tgbot.RenderActionType{
			&tgbot.RenderActionReplace{RenderedElement: rm1, NewElement: m4},
			&tgbot.RenderActionReplace{RenderedElement: rm2, NewElement: m1},
			&tgbot.RenderActionReplace{RenderedElement: rm3, NewElement: m2},
			&tgbot.RenderActionCreate{NewElement: m3},
		})

}

func TestGetRenderActionsBasic(t *testing.T) {
	// r := make([]tgbot.RenderedElement, 0)
	// r = append(r, &tgbot.RenderedUserMessage{})
	// n := make([]tgbot.OutcomingMessageType, 0)

	Check(t,
		[]tgbot.RenderedElement{},
		[]tgbot.OutcomingMessage{m1},
		[]tgbot.RenderActionType{
			&tgbot.RenderActionCreate{NewElement: m1},
		})

	Check(t,
		[]tgbot.RenderedElement{},
		[]tgbot.OutcomingMessage{m1, m2, m3},
		[]tgbot.RenderActionType{
			&tgbot.RenderActionCreate{NewElement: m1},
			&tgbot.RenderActionCreate{NewElement: m2},
			&tgbot.RenderActionCreate{NewElement: m3},
		})

	Check(t,
		[]tgbot.RenderedElement{
			&tgbot.RenderedBotMessage[any]{
				OutcomingTextMessage: m2,
			},
		},
		[]tgbot.OutcomingMessage{m2},
		[]tgbot.RenderActionType{
			&tgbot.RenderActionKeep{
				RenderedElement: &tgbot.RenderedBotMessage[any]{
					OutcomingTextMessage: m2,
				},
				NewElement: m2,
			},
		})

	// rm1 is supposed to be replaced with m2
	Check(t,
		[]tgbot.RenderedElement{rm1},
		[]tgbot.OutcomingMessage{m2},
		[]tgbot.RenderActionType{
			&tgbot.RenderActionReplace{
				RenderedElement: rm1,
				NewElement:      m2,
			},
		})

	// rm1 is supposed to be removed
	Check(t,
		[]tgbot.RenderedElement{rm1},
		[]tgbot.OutcomingMessage{},
		[]tgbot.RenderActionType{
			&tgbot.RenderActionRemove{
				RenderedElement: rm1,
			},
		})

}

type MockRenderer struct {
	OutcomingMessages []tgbot.OutcomingMessage
}

func (mr *MockRenderer) Message(ctx context.Context, props *tgbot.ChatRendererMessageProps) (*models.Message, error) {
	return &models.Message{}, nil
}

func (mr *MockRenderer) Delete(messageId int) error {
	return nil
}

func TestCreate(t *testing.T) {
	re := &MockRenderer{}

	m := tgbot.NewOutcomingTextMessage[int]("message 1")
	b1 := tgbot.Button("button 1", func() int { return 1 }, "button 1", false)
	m.AddButton(b1)

	tgbot.ExecuteRenderActions[int](
		context.Background(),
		re,
		[]tgbot.RenderActionType{
			&tgbot.RenderActionCreate{
				NewElement: m,
			},
		},
	)
}
