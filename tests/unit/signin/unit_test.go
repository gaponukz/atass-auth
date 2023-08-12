package signin

import (
	"auth/src/application/dto"
	"auth/src/application/usecases/signin"
	"auth/tests/unit/mocks"
	"testing"
)

func TestLogin(t *testing.T) {
	db := mocks.NewMockStorage()
	hash := func(s string) string { return s }

	service := signin.NewSigninService(db, hash)

	user, err := db.Create(dto.CreateUserDTO{Gmail: "test@user.ua", Password: "12345", FullName: "Sometest"})
	if err != nil {
		t.Error(err.Error())
	}

	u, err := service.Login("test@user.ua", "12345")
	if err != nil {
		t.Error(err.Error())
	}

	if user.ID != u.ID {
		t.Errorf("expected %s got %s", u.ID, user.ID)
	}

	if user.FullName != u.FullName {
		t.Errorf("expected %s got %s", user.FullName, u.FullName)
	}
}
