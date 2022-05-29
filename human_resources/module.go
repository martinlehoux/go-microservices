package human_resources

import (
	"crypto/rsa"
	"go-microservices/common"
	"go-microservices/human_resources/group"
	"go-microservices/human_resources/user"
)

func Bootstrap(logger common.Logger, publicKey rsa.PublicKey) *HumanResourcesHttpController {
	userStore := user.NewSqlUserStore("human_resources")
	groupStore := group.NewMongoGroupStore("human_resources")
	service := NewHumanResourcesService(&userStore, &groupStore, logger)
	controller := NewHumanResourcesHttpController(&service, publicKey)
	return &controller
}
