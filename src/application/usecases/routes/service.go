package routes

import (
	"auth/src/entities"
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

func (r routesService) AddRoute(userID string, path entities.Path) error {
	user, err := r.db.ByID(userID)
	if err != nil {
		return err
	}

	user.PurchasedRouteIds = append(user.PurchasedRouteIds, path)

	return r.db.Update(user)
}

func (r routesService) DisplayUserRoutes(userID string) ([]entities.Path, error) {
	user, err := r.db.ByID(userID)
	if err != nil {
		return nil, err
	}

	return user.PurchasedRouteIds, nil
}
