package settings

import (
	"auth/src/application/dto"
	"auth/src/application/usecases/settings"
	"auth/tests/unit/mocks"
	"testing"
)

func TestUpdate(t *testing.T) {
	db := mocks.NewMockStorage()
	service := settings.NewSettingsService(db)

	user, err := db.Create(dto.CreateUserDTO{Gmail: "test@gmail.com", Password: "12345", FullName: "Test1"})
	if err != nil {
		t.Fatal(err.Error())
	}

	err = service.UpdateWithFields(user.ID, dto.UpdateUserDTO{
		FullName:            "Test2",
		Phone:               "12345",
		AllowsAdvertisement: true,
	})
	if err != nil {
		t.Error(err.Error())
	}

	u, err := db.ByID(user.ID)
	if err != nil {
		t.Error(err.Error())
	}

	if u.FullName != "Test2" {
		t.Errorf("expected Test2 got %s", u.FullName)
	}

	if u.Phone != "12345" {
		t.Errorf("expected 12345 got %s", u.Phone)
	}

	if u.AllowsAdvertisement != true {
		t.Errorf("expected true got %t", u.AllowsAdvertisement)
	}
}
