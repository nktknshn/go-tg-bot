package todo

import (
	"regexp"
	"strconv"

	tgbot "github.com/nktknshn/go-tg-bot"
)

type PageTodoListState struct {
	CandidateItem string

	Selected      bool
	SelectedIndex int
}

type S = PageTodoListState

type PageTodoList struct {
	Context AppGlobalContext
	State   tgbot.State[PageTodoListState]
}

var rexItemIdex = regexp.MustCompile(`^/(\d+)`)

// Select TodoList from the global context.
// When TodoList is updated, the page will be re-rendered
func (a *PageTodoList) Selector() TodoList {
	return a.Context.TodoList
}

func (a *PageTodoList) Render(o tgbot.OO) {
	tdl := a.Selector()

	// initialize the component state
	state := a.State.Init(PageTodoListState{})
	resetState := state.Set(func(ptls S) S {
		return S{}
	})

	candidateItem := state.Get().CandidateItem
	hasCandidateItem := candidateItem != ""

	isItemSelected := state.Get().Selected
	selectedIndex := state.Get().SelectedIndex

	o.InputHandler(func(s string) any {

		if hasCandidateItem {
			return tgbot.Next{}
		}

		if rexItemIdex.Match([]byte(s)) {
			idxStr := rexItemIdex.FindStringSubmatch(s)[1]
			idx, _ := strconv.Atoi(idxStr)

			return state.Set(func(ptls S) S {
				return S{
					Selected:      true,
					SelectedIndex: idx,
				}
			})
		}

		return state.Set(func(ptls S) S {
			return S{CandidateItem: s}
		})
	})

	o.MessagePart("Todo list:")

	for idx, item := range tdl.Items {
		if item.Done {
			o.MessagePartf("/%v [x] %v", idx, item.Text)
		} else {
			o.MessagePartf("/%v [ ] %v", idx, item.Text)
		}
	}

	if isItemSelected {
		selectedItem := tdl.Items[selectedIndex]
		o.MessagePartf("Selected: %v", selectedItem.Text)

		o.Button("Done", func() any {
			return []any{
				resetState,
				ActionMarkDone{ItemIndex: selectedIndex},
			}
		})

		o.Button("Delete", func() any {
			return ActionItemDelete{ItemIndex: selectedIndex}
		})

		o.Button("Cancel", func() any {
			return resetState
		})

	} else {
		o.MessageComplete()
	}

	if hasCandidateItem {

		o.Messagef("Add %v?", candidateItem)
		o.Button("Yes", func() any {
			return []any{
				resetState,
				ActionAddTodoItem{Text: candidateItem},
			}
		})
		o.Button("No", func() any {
			return resetState
		})

	}
}
