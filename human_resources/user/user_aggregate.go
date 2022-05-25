package user

import (
	"go-microservices/common"
)

type UserID struct{ common.ID }

type User struct {
	id            UserID
	preferredName string
	email         string
}

func New(preferredName string, email string) User {
	return User{
		id:            UserID{common.CreateID()},
		preferredName: preferredName,
		email:         email,
	}
}

func (user *User) Rename(preferredName string) {
	user.preferredName = preferredName
}

func (user *User) GetID() UserID {
	return user.id
}

func (user *User) GetEmail() string {
	return user.email
}
