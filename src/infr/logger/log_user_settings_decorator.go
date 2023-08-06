package logger

import (
	"auth/src/application/dto"
	"fmt"
)

type settingsService interface {
	UpdateWithFields(string, dto.UpdateUserDTO) error
}

type logSettingsService struct {
	s settingsService
	l logger
}

func NewLogSettingsServiceDecorator(s settingsService, l logger) logSettingsService {
	return logSettingsService{s: s, l: l}
}

func (s logSettingsService) UpdateWithFields(id string, fields dto.UpdateUserDTO) error {
	err := s.s.UpdateWithFields(id, fields)
	if err != nil {
		s.l.Error(fmt.Sprintf("Can not update user %s: %v", id, err))
		return err
	}

	s.l.Info(fmt.Sprintf("User %s updated", id))
	return nil
}
