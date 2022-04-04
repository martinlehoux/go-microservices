package user

import (
	"testing"
)

func testModule() (*UserService, UserStore) {
	userStore := bootstrapFakeUserStore()
	service := UserService{userStore: &userStore}
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
	userStore.save(NewUser(NewUserPayload{Email: "test@test.com"}))

	err := service.Register("test@test.com")

	if err == nil {
		t.Errorf("expected Register to error but got nothing")
	}
}
