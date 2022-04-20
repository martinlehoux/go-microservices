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

type NewUserPayload struct {
	PreferredName string
	Email         string
}

func NewUser(payload NewUserPayload) User {
	return User{
		id:            UserID{common.CreateID()},
		preferredName: payload.PreferredName,
		email:         payload.Email,
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
