package logger

import (
	"auth/src/entities"
	"fmt"
)

type routesService interface {
	AddRoute(string, entities.Path) error
}

type logger interface {
	Error(string)
	Info(string)
}

type logRoutesService struct {
	s routesService
	l logger
}

func NewLogAddRouteDecorator(s routesService, l logger) logRoutesService {
	return logRoutesService{s: s, l: l}
}

func (s logRoutesService) AddRoute(userID string, path entities.Path) error {
	err := s.s.AddRoute(userID, path)
	if err != nil {
		s.l.Error(fmt.Sprintf("Can not add route %s to user %s: %v", path.RootRouteID, userID, err))
	}

	s.l.Info(fmt.Sprintf("Add route %s to user %s", path.RootRouteID, userID))
	return nil
}
