package todo

type TodoItem struct {
	Text string
	Done bool
	Tags []string
}

type TodoList struct {
	Items []TodoItem
}
