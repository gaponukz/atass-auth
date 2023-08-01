package logger

import "fmt"

type routesService interface {
	AddRoute(userID, routeID string) error
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

func (s logRoutesService) AddRoute(userID, routeID string) error {
	err := s.s.AddRoute(userID, routeID)
	if err != nil {
		s.l.Error(fmt.Sprintf("Can not add route %s to user %s: %v", routeID, userID, err))
	}

	s.l.Info(fmt.Sprintf("Add route %s to user %s", routeID, userID))
	return nil
}
