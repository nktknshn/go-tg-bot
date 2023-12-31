package todo

import (
	"regexp"
	"strconv"

	tgbot "github.com/nktknshn/go-tg-bot/tgbot"
	"github.com/nktknshn/go-tg-bot/tgbot/component"
)

var rexItemIdex = regexp.MustCompile(`^/(\d+)`)

type PageTodoListState struct {
	CandidateItem string

	Selected      bool
	SelectedIndex int
}

type s = PageTodoListState

type PageTodoList struct {
	Context TodoGlobalContext
	State   component.State[PageTodoListState]
}

// Select TodoList from the global context.
// When TodoList is updated, the page will be re-rendered
func (a *PageTodoList) Selector() TodoList {
	return a.Context.TodoList
}

func (a *PageTodoList) Render(o tgbot.O) {
	tdl := a.Selector()

	// initialize the component state
	state := a.State.Init(PageTodoListState{})
	resetState := state.Set(func(ptls s) s {
		return s{}
	})

	candidateItem := state.Get().CandidateItem
	hasCandidateItem := candidateItem != ""

	isItemSelected := state.Get().Selected
	selectedIndex := state.Get().SelectedIndex

	o.InputHandler(func(str string) any {

		if hasCandidateItem {
			return tgbot.Next()
		}

		if rexItemIdex.Match([]byte(str)) {
			idxStr := rexItemIdex.FindStringSubmatch(str)[1]
			idx, _ := strconv.Atoi(idxStr)

			return state.Set(func(ptls s) s {
				return s{
					Selected:      true,
					SelectedIndex: idx,
				}
			})
		}

		return state.Set(func(ptls s) s {
			return s{CandidateItem: str}
		})
	})

	if tdl.IsEmpty() && !hasCandidateItem {
		o.Message("Send first todo item:")
		return
	}

	if !tdl.IsEmpty() {
		for idx, item := range tdl.Items {
			if item.Done {
				o.MessagePartf("/%v 🟢 %v", idx, item.Text)
			} else {
				o.MessagePartf("/%v ⭕️ %v", idx, item.Text)
			}
		}

	}

	if isItemSelected {
		selectedItem := tdl.Items[selectedIndex]
		o.MessagePartf("\nSelected: %v", selectedItem.Text)

		var btnCaption string

		if selectedItem.Done {
			btnCaption = "⭕️ Undone"
		} else {
			btnCaption = "🟢 Done"
		}

		o.Button(btnCaption, func() any {
			return []any{
				resetState,
				ActionToggle{ItemIndex: selectedIndex},
			}
		})

		o.Button("Delete", func() any {
			return ActionItemDelete{ItemIndex: selectedIndex}
		})

		o.Button("Cancel", func() any {
			return resetState
		})

	}

	o.MessageComplete()

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
