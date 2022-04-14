package human_resources

import "go-microservices/human_resources/user"

func Bootstrap() *HumanResourcesService {
	store := user.NewSqlUserStore()
	return &HumanResourcesService{
		userStore: &store,
	}
}
