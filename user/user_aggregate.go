package user

import "go-microservices/common"

type UserID common.ID

type User struct {
	Id            UserID
	PreferredName string
	Email         string
}

type NewUserPayload struct {
	PreferredName string
	Email         string
}

func NewUser(payload NewUserPayload) User {
	return User{
		Id:            UserID(common.CreateID()),
		PreferredName: payload.PreferredName,
		Email:         payload.Email,
	}
}
