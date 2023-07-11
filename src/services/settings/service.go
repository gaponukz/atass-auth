package settings

import (
	"auth/src/dto"
	"auth/src/entities"
)

type storage interface {
	ByID(string) (entities.UserEntity, error)
	Update(entities.UserEntity) error
}

type settingsService struct {
	db storage
}

func NewSettingsService(db storage) settingsService {
	return settingsService{db: db}
}

func (s settingsService) UpdateWithFields(id string, fields dto.UpdateUserDTO) error {
	user, err := s.db.ByID(id)
	if err != nil {
		return err
	}

	user.FullName = fields.FullName
	user.Phone = fields.Phone
	user.AllowsAdvertisement = fields.AllowsAdvertisement

	return s.db.Update(user)
}

func (s settingsService) SubscribeUserToRoutes(id string, routeID string) error {
	user, err := s.db.ByID(id)
	if err != nil {
		return err
	}

	user.PurchasedRouteIds = append(user.PurchasedRouteIds, routeID)

	return s.db.Update(user)
}
