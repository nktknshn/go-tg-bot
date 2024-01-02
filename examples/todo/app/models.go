package todo

import "github.com/gotd/td/tg"

type TodoItem struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
	// Tags []string `json:"tags"`
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

type User struct {
	ChatID     int64 `json:"chat_id"`
	AccessHash int64 `json:"access_hash"`

	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	TodoList TodoList `json:"todo_list"`
}

func UserFromTgUser(tgUser *tg.User) *User {
	return &User{
		ChatID:     tgUser.ID,
		AccessHash: tgUser.AccessHash,
		Username:   tgUser.Username,
		FirstName:  tgUser.FirstName,
		LastName:   tgUser.LastName,
	}
}
