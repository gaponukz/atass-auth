package routes

import (
	"auth/src/entities"
	"auth/src/errors"
)

type storage interface {
	ByID(string) (entities.UserEntity, error)
	Update(entities.UserEntity) error
}

type routesService struct {
	db storage
}

func NewRoutesService(db storage) routesService {
	return routesService{db: db}
}

func (s routesService) SubscribeUserToRoutes(userID, routeID string) error {
	user, err := s.db.ByID(userID)
	if err != nil {
		return err
	}

	user.PurchasedRouteIds = append(user.PurchasedRouteIds, routeID)

	return s.db.Update(user)
}

func (s routesService) UnsubscribeUserFromRoutes(userID, routeID string) error {
	user, err := s.db.ByID(userID)
	if err != nil {
		return err
	}

	indexToRemove := -1
	for i, purchasedRouteID := range user.PurchasedRouteIds {
		if purchasedRouteID == routeID {
			indexToRemove = i
			break
		}
	}

	if indexToRemove == -1 {
		return errors.ErrRouteNotFound
	}

	user.PurchasedRouteIds = append(user.PurchasedRouteIds[:indexToRemove], user.PurchasedRouteIds[indexToRemove+1:]...)

	return s.db.Update(user)
}

func (s routesService) GetRoutes(userID string) ([]string, error) {
	user, err := s.db.ByID(userID)
	if err != nil {
		return nil, err
	}

	return user.PurchasedRouteIds, nil
}
