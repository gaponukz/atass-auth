package logger

import (
	"auth/src/domain/entities"
	"auth/src/domain/errors"
	"fmt"
)

type signinService interface {
	Login(string, string) (entities.User, error)
}

type logSigninService struct {
	s signinService
	l logger
}

func NewLogSigninServiceDecorator(s signinService, l logger) logSigninService {
	return logSigninService{s: s, l: l}
}

func (s logSigninService) Login(gmail string, password string) (entities.User, error) {
	user, err := s.s.Login(gmail, password)
	if err != nil {
		if err == errors.ErrUserNotFound {
			return user, err
		}

		s.l.Error(fmt.Sprintf("Can not login to %s: %v", gmail, err))
		return user, err
	}

	return user, nil
}
