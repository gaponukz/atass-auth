package logger

import (
	"auth/src/dto"
	"auth/src/errors"
	"fmt"
)

type resetPasswordService interface {
	NotifyUser(string) (string, error)
	AddUserToTemporaryStorage(dto.GmailWithKeyPairDTO) error
	CancelPasswordResetting(dto.GmailWithKeyPairDTO) error
	ChangeUserPassword(dto.PasswordResetDTO) error
}

type logResetPasswordService struct {
	s resetPasswordService
	l logger
}

func NewLogResetPasswordServiceDecorator(s resetPasswordService, l logger) logResetPasswordService {
	return logResetPasswordService{s: s, l: l}
}

func (s logResetPasswordService) NotifyUser(gmail string) (string, error) {
	code, err := s.s.NotifyUser(gmail)
	if err != nil {
		s.l.Error(fmt.Sprintf("Can not send code to %s: %v", gmail, err))
	}

	return code, err
}

func (s logResetPasswordService) AddUserToTemporaryStorage(d dto.GmailWithKeyPairDTO) error {
	err := s.s.AddUserToTemporaryStorage(d)
	if err != nil {
		if err == errors.ErrUserNotFound {
			return err
		}

		s.l.Error(fmt.Sprintf("Can not AddUserToTemporaryStorage %s: %v", d.Gmail, err))
	}

	return err
}

func (s logResetPasswordService) CancelPasswordResetting(d dto.GmailWithKeyPairDTO) error {
	err := s.s.CancelPasswordResetting(d)
	if err != nil {
		if err == errors.ErrPasswordResetRequestMissing {
			return err
		}

		s.l.Error(fmt.Sprintf("Can not CancelPasswordResetting %s: %v", d.Gmail, err))
	}

	return err
}

func (s logResetPasswordService) ChangeUserPassword(d dto.PasswordResetDTO) error {
	err := s.s.ChangeUserPassword(d)
	if err != nil {
		if err == errors.ErrPasswordResetRequestMissing {
			return err
		}

		s.l.Error(fmt.Sprintf("Can not ChangeUserPassword %s: %v", d.Gmail, err))
	}

	return err
}
