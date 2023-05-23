package models

import "time"

type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

func GetUserByID(id uint) (*User, error) {
	return nil, nil
}

func SaveUser(user *User) error {
	return nil
}
