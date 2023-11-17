package todo

import tgbot "github.com/nktknshn/go-tg-bot"

type PageTodoListState struct {
	CandidateItem string
}

type S = PageTodoListState

type PageTodoList struct {
	Context AppGlobalContext
	State   tgbot.State[PageTodoListState]
}

// Select TodoList from the global context.
// When TodoList is updated, the page will be re-rendered
func (a *PageTodoList) Selector() TodoList {
	return a.Context.TodoList
}

func (a *PageTodoList) Render(o tgbot.OO) {
	tdl := a.Selector()

	// initialize the component state
	state := a.State.Init(PageTodoListState{})
	resetCandiate := state.Set(func(ptls S) S {
		return S{}
	})
	candidateItem := state.Get().CandidateItem
	hasCandidateItem := candidateItem != ""

	o.InputHandler(func(s string) any {

		if hasCandidateItem {
			return tgbot.Next{}
		}

		return state.Set(func(ptls S) S {
			return S{CandidateItem: s}
		})
	})

	o.MessagePart("Todo list:")

	for idx, item := range tdl.Items {
		o.MessagePartf("/%v %v", idx, item.Text)
	}

	o.MessageComplete()

	if hasCandidateItem {

		o.Messagef("Add %v?", candidateItem)
		o.Button("Yes", func() any {
			return []any{
				resetCandiate,
				ActionAddTodoItem{Text: candidateItem},
			}
		})
		o.Button("No", func() any {
			return resetCandiate
		})

	}
}
