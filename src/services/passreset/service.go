package passreset

import (
	"auth/src/dto"
	"auth/src/entities"
	"auth/src/errors"
	"auth/src/utils"
)

type updateAndReadAbleStorage interface {
	ReadAll() ([]entities.UserEntity, error)
	Update(entities.UserEntity) error
}

type gmailKeyPairStorage interface {
	Create(dto.GmailWithKeyPairDTO) error
	Delete(dto.GmailWithKeyPairDTO) error
	GetByUniqueKey(string) (dto.GmailWithKeyPairDTO, error)
}

type resetPasswordService struct {
	temporaryStorage gmailKeyPairStorage
	userStorage      updateAndReadAbleStorage
	notify           func(gmail, key string) error
	generateCode     func() string
	hash             func(string) string
}

func NewResetPasswordService(
	userStorage updateAndReadAbleStorage,
	temporaryStorage gmailKeyPairStorage,
	notify func(gmail, key string) error,
	hash func(string) string,
	generateCode func() string,
) *resetPasswordService {
	return &resetPasswordService{
		temporaryStorage: temporaryStorage,
		userStorage:      userStorage,
		notify:           notify,
		generateCode:     generateCode,
		hash:             hash,
	}
}

func (s resetPasswordService) NotifyUser(userGmail string) (string, error) {
	key := s.generateCode()
	err := s.notify(userGmail, key)

	return key, err
}

func (s resetPasswordService) AddUserToTemporaryStorage(user dto.GmailWithKeyPairDTO) error {
	users, err := s.userStorage.ReadAll()
	if err != nil {
		return err
	}

	isGmailExist := utils.IsExist(users, func(u entities.UserEntity) bool {
		return u.Gmail == user.Gmail
	})

	if !isGmailExist {
		return errors.ErrUserNotFound
	}

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

	userToUpdate, err := utils.Find(users, func(u entities.UserEntity) bool {
		return u.Gmail == data.Gmail
	})

	if err != nil {
		return err
	}

	userToUpdate.Password = s.hash(data.Password)

	return s.userStorage.Update(userToUpdate)
}
