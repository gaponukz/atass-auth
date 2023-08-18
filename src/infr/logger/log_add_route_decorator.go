package logger

import (
	"auth/src/domain/entities"
	"fmt"
)

type addRouteService interface {
	AddRoute(string, entities.Path) error
}

type logAddRouteService struct {
	s addRouteService
	l logger
}

func NewLogAddRouteDecorator(s addRouteService, l logger) logAddRouteService {
	return logAddRouteService{s: s, l: l}
}

func (s logAddRouteService) AddRoute(userID string, path entities.Path) error {
	err := s.s.AddRoute(userID, path)
	if err != nil {
		s.l.Error(fmt.Sprintf("Can not add route %s to user %s: %v", path.RootRouteID, userID, err))
	}

	s.l.Info(fmt.Sprintf("Add route %s to user %s", path.RootRouteID, userID))
	return nil
}
