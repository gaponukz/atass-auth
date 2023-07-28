package routes

import (
	"auth/src/entities"
	"auth/src/errors"
)

type db interface {
	ByID(string) (entities.UserEntity, error)
	Update(entities.UserEntity) error
}

type routesService struct {
	db db
}

func NewRoutesService(db db) routesService {
	return routesService{db: db}
}

func (r routesService) AddRoute(userID, routeID string) error {
	user, err := r.db.ByID(userID)
	if err != nil {
		return err
	}

	user.PurchasedRouteIds = append(user.PurchasedRouteIds, routeID)

	return r.db.Update(user)
}

func (r routesService) DeleteRoute(userID, routeID string) error {
	user, err := r.db.ByID(userID)
	if err != nil {
		return err
	}

	index := -1
	for i, id := range user.PurchasedRouteIds {
		if id == routeID {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.ErrRouteNotFound
	}

	user.PurchasedRouteIds = append(user.PurchasedRouteIds[:index], user.PurchasedRouteIds[index+1:]...)

	return r.db.Update(user)
}

func (r routesService) DisplayUserRoutes(userID string) ([]string, error) {
	user, err := r.db.ByID(userID)
	if err != nil {
		return nil, err
	}

	return user.PurchasedRouteIds, nil
}
