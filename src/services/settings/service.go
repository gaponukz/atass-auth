package settings

import "auth/src/entities"

type storage interface {
	Update(entities.UserEntity) error
}

type settingsService struct {
	db storage
}

func NewSettingsService(db storage) settingsService {
	return settingsService{db: db}
}

func (s settingsService) Update(user entities.UserEntity) error {
	return s.db.Update(user)
}
