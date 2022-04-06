package authentication

import "testing"

func testModule() *AuthenticationService {
	accountStore := bootstrapFakeAccountStore()
	service := AuthenticationService{accountStore: &accountStore}
	return &service
}

func TestRegisterNewAccount(t *testing.T) {
	service := testModule()

	err := service.Register("identifier", "password")

	if err != nil {
		t.Errorf("expected Register not to error but got: %s", err)
	}

	_, err = service.accountStore.loadForIdentifier("identifier")
	if err != nil {
		t.Errorf("expected account to be saved but got: %s", err)
	}
}

func TestRegisterExistingAccount(t *testing.T) {
	service := testModule()
	service.accountStore.save(NewAccount("identifier", []byte("password")))

	err := service.Register("identifier", "password")

	if err == nil {
		t.Errorf("expected Register to error but got nothing")
	}
	if err.Error() != "identifier already used" {
		t.Errorf("expected error to be identifier already used but got: %s", err)
	}
}
