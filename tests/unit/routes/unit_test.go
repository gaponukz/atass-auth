package routes

import (
	"auth/src/application/dto"
	"auth/src/application/usecases/routes"
	"auth/src/domain/entities"
	"auth/src/domain/errors"
	"auth/tests/unit/mocks"
	"testing"
)

func TestAddRoute(t *testing.T) {
	db := mocks.NewMockStorage()
	service := routes.NewAddRouteService(db)

	err := service.AddRoute("1", entities.Path{RootRouteID: "123", MoveFromID: "1", MoveToID: "2"})
	if err != errors.ErrUserNotFound {
		t.Errorf("can add route for not existens user")
	}

	user, err := db.Create(dto.CreateUserDTO{})
	if err != nil {
		t.Error(err.Error())
	}

	err = service.AddRoute(user.ID, entities.Path{RootRouteID: "123", MoveFromID: "1", MoveToID: "2"})
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
		if user.PurchasedRouteIds[0].RootRouteID != "123" {
			t.Errorf("Can not add route to user")
		}
	}
}
