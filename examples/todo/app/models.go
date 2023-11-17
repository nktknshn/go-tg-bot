package todo

type TodoItem struct {
	Text string
	Done bool
	Tags []string
}

type TodoList struct {
	Items []TodoItem
}

func (tdl TodoList) Count() int {
	return len(tdl.Items)
}

// is empty
func (tdl TodoList) IsEmpty() bool {
	return tdl.Count() == 0
}
