package human_resources

import (
	"crypto/rsa"
	"go-microservices/human_resources/group"
	"go-microservices/human_resources/user"
)

func Bootstrap(rootPath string, publicKey rsa.PublicKey) *HumanResourcesHttpController {
	userStore := user.NewSqlUserStore()
	groupStore := group.NewMongoGroupStore()
	service := NewHumanResourcesService(&userStore, &groupStore)
	controller := HumanResourcesHttpController{
		humanResourcesService: service,
		publicKey:             publicKey,
		rootPath:              rootPath,
	}
	return &controller
}
