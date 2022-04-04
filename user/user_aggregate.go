package user

import "go-microservices/common"

type User struct {
	Id            common.ID
	PreferredName string
	Email         string
}

type NewUserPayload struct {
	PreferredName string
	Email         string
}

func NewUser(payload NewUserPayload) User {
	return User{
		Id:            common.CreateID(),
		PreferredName: payload.PreferredName,
		Email:         payload.Email,
	}
}
