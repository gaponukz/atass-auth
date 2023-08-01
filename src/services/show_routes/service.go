package show_routes

import "auth/src/entities"

type db interface {
	ByID(string) (entities.UserEntity, error)
}

type service struct {
	db db
}

func NewShowRoutesService(db db) service {
	return service{db: db}
}

func (s service) ShowRoutes(id string) ([]string, error) {
	user, err := s.db.ByID(id)
	if err != nil {
		return nil, err
	}

	return user.PurchasedRouteIds, nil
}
