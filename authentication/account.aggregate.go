package authentication

import "go-microservices/common"

type Account struct {
	Id       common.ID
	TargetId common.ID
}

func NewAccount(targetId common.ID) Account {
	return Account{
		Id:       common.CreateID(),
		TargetId: targetId,
	}
}
