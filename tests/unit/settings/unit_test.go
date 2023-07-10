package settings

import (
	"auth/src/entities"
	"auth/src/services/settings"
	"auth/tests/unit/mocks"
	"testing"
)

func TestUpdate(t *testing.T) {
	db := mocks.NewMockStorage()
	service := settings.NewSettingsService(db)

	user, err := db.Create(entities.User{Gmail: "test@gmail.com", Password: "12345", FullName: "Test1"})
	if err != nil {
		t.Fatal(err.Error())
	}

	user.FullName = "Test2"
	err = service.Update(user)
	if err != nil {
		t.Error(err.Error())
	}

	u, err := db.ByID(user.ID)
	if err != nil {
		t.Error(err.Error())
	}

	if u.FullName != user.FullName {
		t.Errorf("expected %s got %s", user.FullName, u.FullName)
	}
}
