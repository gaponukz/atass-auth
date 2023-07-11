package settings

import (
	"auth/src/entities"
	"errors"
)

type storage interface {
	ByID(string) (entities.UserEntity, error)
	Update(entities.UserEntity) error
}

type settingsService struct {
	db storage
}

type updateUserDTO struct {
	FullName            string `json:"fullName"`
	Phone               string `json:"phone"`
	AllowsAdvertisement bool   `json:"allowsAdvertisement"`
}

func NewSettingsService(db storage) settingsService {
	return settingsService{db: db}
}

func (s settingsService) UpdateWithFields(id string, fields interface{}) error {
	dto, ok := fields.(updateUserDTO)
	if !ok {
		return errors.New("expected updateUserDTO")
	}

	user, err := s.db.ByID(id)
	if err != nil {
		return err
	}

	user.FullName = dto.FullName
	user.Phone = dto.Phone
	user.AllowsAdvertisement = dto.AllowsAdvertisement

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
