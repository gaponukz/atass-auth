package resetPassword

import (
	"auth/src/entities"
	"fmt"
)

type updatePasswordGetByGmailAbleStorage interface {
	GetByGmail(string) (entities.User, error)
	UpdatePassword(entities.User, string) error
}

type gmailKeyPairStorage interface {
	Create(entities.GmailWithKeyPair) error
	Delete(entities.GmailWithKeyPair) error
	GetByUniqueKey(string) (entities.GmailWithKeyPair, error)
}

type resetPasswordService struct {
	temporaryStorage gmailKeyPairStorage
	userStorage      updatePasswordGetByGmailAbleStorage
	notify           func(gmail, key string) error
	generateCode     func() string
}

func NewResetPasswordService(
	userStorage updatePasswordGetByGmailAbleStorage,
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
	_, err := s.userStorage.GetByGmail(user.Gmail)

	if err != nil {
		return err
	}

	return s.temporaryStorage.Create(user)
}

func (s resetPasswordService) ChangeUserPassword(user entities.GmailWithKeyPair, newPassword string) error {
	_, err := s.temporaryStorage.GetByUniqueKey(user.Key)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	err = s.temporaryStorage.Delete(user)
	if err != nil {
		return fmt.Errorf("could not remove user")
	}

	realUser, err := s.userStorage.GetByGmail(user.Gmail)
	if err != nil {
		return err
	}

	return s.userStorage.UpdatePassword(realUser, newPassword)
}
