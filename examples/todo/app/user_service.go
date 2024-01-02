package todo

import (
	"encoding/json"
	"os"
	"sync"
)

type UserService interface {
	GetUser(chatID int64) (*User, error)
	SaveUser(user *User) error
}

type UserServiceJson struct {
	users map[int64]*User
	file  string

	lock *sync.Mutex
}

func NewUserServiceJson(file string) *UserServiceJson {
	return &UserServiceJson{
		users: make(map[int64]*User),
		file:  file,
		lock:  &sync.Mutex{},
	}
}

func (s *UserServiceJson) GetUser(chatID int64) (*User, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	user, ok := s.users[chatID]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (s *UserServiceJson) SaveUser(user *User) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.users[user.ChatID] = user

	err := s.save()

	return err
}

func (s *UserServiceJson) json() ([]byte, error) {
	return json.Marshal(s.users)
}

func (s *UserServiceJson) save() error {
	// save to file

	j, err := s.json()

	if err != nil {
		return err
	}

	err = os.WriteFile(s.file, j, 0644)

	return err
}

func (s *UserServiceJson) Load() error {
	// load from file

	j, err := os.ReadFile(s.file)

	if err != nil {
		return err
	}

	err = json.Unmarshal(j, &s.users)

	if err != nil {
		return err
	}

	return nil
}
