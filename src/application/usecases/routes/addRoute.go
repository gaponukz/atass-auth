package routes

import "auth/src/domain/entities"

type db interface {
	ByID(string) (entities.User, error)
	Update(entities.User) error
}

type addRouteService struct {
	db db
}

func NewAddRouteService(db db) addRouteService {
	return addRouteService{db: db}
}

func (r addRouteService) AddRoute(userID string, path entities.Path) error {
	user, err := r.db.ByID(userID)
	if err != nil {
		return err
	}

	user.PurchasedRouteIds = append(user.PurchasedRouteIds, path)

	return r.db.Update(user)
}
