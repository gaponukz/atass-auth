package logger

import (
	"auth/src/application/dto"
	"auth/src/domain/errors"
	"fmt"
)

type signupService interface {
	SendGeneratedCode(string) (string, error)
	AddUserToTemporaryStorage(dto.GmailWithKeyPairDTO) error
	RegisterUserOnRightCode(dto.SignUpDTO) (string, error)
}

type logSignupService struct {
	s signupService
	l logger
}

func NewLogSignupServiceDecorator(s signupService, l logger) logSignupService {
	return logSignupService{s: s, l: l}
}

func (s logSignupService) SendGeneratedCode(gmail string) (string, error) {
	code, err := s.s.SendGeneratedCode(gmail)
	if err != nil {
		s.l.Error(fmt.Sprintf("Can not send generated code to %s: %v", gmail, err))
	}

	return code, err
}

func (s logSignupService) AddUserToTemporaryStorage(d dto.GmailWithKeyPairDTO) error {
	err := s.s.AddUserToTemporaryStorage(d)
	if err != nil {
		if err == errors.ErrUserAlreadyExists {
			return err
		}

		s.l.Error(fmt.Sprintf("Can not add %s to temporary storage: %v", d.Gmail, err))
		return err
	}

	return nil
}

func (s logSignupService) RegisterUserOnRightCode(d dto.SignUpDTO) (string, error) {
	id, err := s.s.RegisterUserOnRightCode(d)
	if err != nil {
		if err == errors.ErrRegisterRequestMissing {
			return id, err
		}
		if err == errors.ErrUserNotValid {
			return id, err
		}

		s.l.Error(fmt.Sprintf("Can not register %s with right code: %v", d.Gmail, err))
		return id, err
	}

	return id, nil
}
