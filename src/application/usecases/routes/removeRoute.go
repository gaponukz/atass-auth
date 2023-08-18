package routes

import (
	"auth/src/domain/entities"
	"auth/src/domain/errors"
)

type deleteRouteService struct {
	db db
}

func NewDeleteRouteService(db db) deleteRouteService {
	return deleteRouteService{db: db}
}

func (r deleteRouteService) DeleteRoute(userID string, path entities.Path) error {
	user, err := r.db.ByID(userID)
	if err != nil {
		return err
	}

	index := -1
	for i, purchasedPath := range user.PurchasedRouteIds {
		if purchasedPath == path {
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
