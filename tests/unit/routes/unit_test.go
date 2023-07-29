package routes

import (
	"auth/src/entities"
	"auth/src/errors"
	"auth/src/services/routes"
	"auth/tests/unit/mocks"
	"testing"
)

func TestAddRoute(t *testing.T) {
	db := mocks.NewMockStorage()
	service := routes.NewRoutesService(db)

	err := service.AddRoute("1", "123")
	if err != errors.ErrUserNotFound {
		t.Errorf("can add route for not existens user")
	}

	user, err := db.Create(entities.User{})
	if err != nil {
		t.Error(err.Error())
	}

	err = service.AddRoute(user.ID, "123")
	if err != nil {
		t.Error(err.Error())
	}

	user, err = db.ByID(user.ID)
	if err != nil {
		t.Error(err.Error())
	}

	if len(user.PurchasedRouteIds) == 0 {
		t.Errorf("Can not add route to user")
	} else {
		if user.PurchasedRouteIds[0] != "123" {
			t.Errorf("Can not add route to user")
		}
	}
}
