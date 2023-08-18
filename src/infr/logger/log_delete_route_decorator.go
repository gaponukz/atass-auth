package logger

import (
	"auth/src/domain/entities"
	"fmt"
)

type deleteRouteService interface {
	DeleteRoute(string, entities.Path) error
}

type logger interface {
	Error(string)
	Info(string)
}

type logRoutesService struct {
	s deleteRouteService
	l logger
}

func NewLogDeleteRouteDecorator(s deleteRouteService, l logger) logRoutesService {
	return logRoutesService{s: s, l: l}
}

func (s logRoutesService) DeleteRoute(userID string, path entities.Path) error {
	err := s.s.DeleteRoute(userID, path)
	if err != nil {
		s.l.Error(fmt.Sprintf("Can not delete route %s to user %s: %v", path.RootRouteID, userID, err))
	}

	s.l.Info(fmt.Sprintf("Delete route %s from user %s", path.RootRouteID, userID))
	return nil
}
