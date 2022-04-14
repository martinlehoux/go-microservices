package human_resources

import (
	"go-microservices/human_resources/user"
	"testing"
)

func testModule() (*HumanResourcesService, user.UserStore) {
	userStore := user.NewFakeUserStore()
	service := HumanResourcesService{userStore: &userStore}
	return &service, &userStore
}

func TestUserRegisterMissingUser(t *testing.T) {
	service, _ := testModule()

	err := service.Register("test@test.com")

	if err != nil {
		t.Errorf("expected Register not to error but got: %s", err)
	}
}

func TestUserRegisterExistingUser(t *testing.T) {
	service, userStore := testModule()
	userStore.Save(user.NewUser(user.NewUserPayload{Email: "test@test.com"}))

	err := service.Register("test@test.com")

	if err == nil {
		t.Errorf("expected Register to error but got nothing")
	}
}
