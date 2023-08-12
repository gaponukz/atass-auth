package settings

import (
	"auth/src/application/dto"
	"auth/src/domain/entities"
)

type storage interface {
	ByID(string) (entities.User, error)
	Update(entities.User) error
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
