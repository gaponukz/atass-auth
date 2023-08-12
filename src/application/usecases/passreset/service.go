package passreset

import (
	"auth/src/application/dto"
	"auth/src/domain/entities"
	"auth/src/domain/errors"
	"auth/src/utils"
)

type notifier interface {
	NotifyUser(entities.User, string) error
}

type updateAndReadAbleStorage interface {
	ReadAll() ([]entities.User, error)
	Update(entities.User) error
}

type gmailKeyPairStorage interface {
	Create(dto.GmailWithKeyPairDTO) error
	Delete(dto.GmailWithKeyPairDTO) error
	GetByUniqueKey(string) (dto.GmailWithKeyPairDTO, error)
}

type resetPasswordService struct {
	temporaryStorage gmailKeyPairStorage
	userStorage      updateAndReadAbleStorage
	notifier         notifier
	generateCode     func() string
	hash             func(string) string
}

func NewResetPasswordService(
	userStorage updateAndReadAbleStorage,
	temporaryStorage gmailKeyPairStorage,
	notifier notifier,
	hash func(string) string,
	generateCode func() string,
) *resetPasswordService {
	return &resetPasswordService{
		temporaryStorage: temporaryStorage,
		userStorage:      userStorage,
		notifier:         notifier,
		generateCode:     generateCode,
		hash:             hash,
	}
}

func (s resetPasswordService) NotifyUser(userGmail string) (string, error) {
	users, err := s.userStorage.ReadAll()
	if err != nil {
		return "", err
	}

	user, err := utils.Find(users, func(u entities.User) bool {
		return u.Gmail == userGmail
	})
	if err != nil {
		return "", errors.ErrUserNotFound
	}

	key := s.generateCode()
	err = s.notifier.NotifyUser(user, key)

	return key, err
}

func (s resetPasswordService) AddUserToTemporaryStorage(user dto.GmailWithKeyPairDTO) error {
	return s.temporaryStorage.Create(user)
}

func (s resetPasswordService) CancelPasswordResetting(user dto.GmailWithKeyPairDTO) error {
	err := s.temporaryStorage.Delete(user)
	if err != nil {
		return errors.ErrPasswordResetRequestMissing
	}

	return nil
}

func (s resetPasswordService) ChangeUserPassword(data dto.PasswordResetDTO) error {
	err := s.temporaryStorage.Delete(dto.GmailWithKeyPairDTO{Gmail: data.Gmail, Key: data.Key})
	if err != nil {
		return errors.ErrPasswordResetRequestMissing
	}

	users, err := s.userStorage.ReadAll()
	if err != nil {
		return err
	}

	userToUpdate, err := utils.Find(users, func(u entities.User) bool {
		return u.Gmail == data.Gmail
	})

	if err != nil {
		return err
	}

	userToUpdate.Password = s.hash(data.Password)

	return s.userStorage.Update(userToUpdate)
}
