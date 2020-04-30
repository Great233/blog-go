package services

import "blog/models"

type User struct {
	Id        uint
	Username  string
	LastLogin int
}

func (u *User) GetByUsername() (*models.User, error) {
	return models.GetUserByUsername(u.Username)
}
