package password_reseting

import (
	"auth/src/entities"
	"auth/src/errors"
	"auth/src/security"
	"auth/src/storage"
	"fmt"
)

type updateAndReadAbleStorage interface {
	ReadAll() ([]entities.UserEntity, error)
	Update(entities.UserEntity) error
}

type gmailKeyPairStorage interface {
	Create(entities.GmailWithKeyPair) error
	Delete(entities.GmailWithKeyPair) error
	GetByUniqueKey(string) (entities.GmailWithKeyPair, error)
}

type resetPasswordService struct {
	temporaryStorage gmailKeyPairStorage
	userStorage      updateAndReadAbleStorage
	notify           func(gmail, key string) error
	generateCode     func() string
}

func NewResetPasswordService(
	userStorage updateAndReadAbleStorage,
	temporaryStorage gmailKeyPairStorage,
	notify func(gmail, key string) error,
	generateCode func() string,
) *resetPasswordService {
	return &resetPasswordService{
		temporaryStorage: temporaryStorage,
		userStorage:      userStorage,
		notify:           notify,
		generateCode:     generateCode,
	}
}

func (s resetPasswordService) NotifyUser(userGmail string) (string, error) {
	key := s.generateCode()
	err := s.notify(userGmail, key)

	return key, err
}

func (s resetPasswordService) AddUserToTemporaryStorage(user entities.GmailWithKeyPair) error {
	users, err := s.userStorage.ReadAll()
	if err != nil {
		return err
	}

	isGmailExist := storage.IsExist(users, func(u entities.UserEntity) bool {
		return u.Gmail == user.Gmail
	})

	if !isGmailExist {
		return fmt.Errorf("gmail %s not found", user.Gmail)
	}

	return s.temporaryStorage.Create(user)
}

func (s resetPasswordService) CancelPasswordResetting(user entities.GmailWithKeyPair) error {
	err := s.temporaryStorage.Delete(user)
	if err != nil {
		return errors.ErrRegisterRequestMissing
	}

	return nil
}

func (s resetPasswordService) ChangeUserPassword(user entities.GmailWithKeyPair, newPassword string) error {
	err := s.temporaryStorage.Delete(user)
	if err != nil {
		return errors.ErrRegisterRequestMissing
	}

	users, err := s.userStorage.ReadAll()
	if err != nil {
		return err
	}

	userToUpdate, err := storage.Find(users, func(u entities.UserEntity) bool {
		return u.Gmail == user.Gmail
	})

	if err != nil {
		return err
	}

	userToUpdate.Password = security.GetSha256(newPassword)

	return s.userStorage.Update(userToUpdate)
}
