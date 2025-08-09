package user

import (
	"time"
)

type Task struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsDone      bool   `json:"isDone"`
	Date        string `json:"date"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Tasks    []Task `json:"tasks,omitempty"`
}

func (u *User) CreateUser(name, password string) User {
	return User{
		ID:       time.Now().Nanosecond(),
		Name:     name,
		Password: password,
	}
}
