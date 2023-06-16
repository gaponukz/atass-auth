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

type ResetPasswordService struct {
	TemporaryStorage gmailKeyPairStorage
	UserStorage      updatePasswordGetByGmailAbleStorage
	Notify           func(gmail, key string) error
	GenerateCode     func() string
}

func (service *ResetPasswordService) NotifyUser(userGmail string) (string, error) {
	key := service.GenerateCode()
	err := service.Notify(userGmail, key)

	return key, err
}

func (service *ResetPasswordService) AddUserToTemporaryStorage(user entities.GmailWithKeyPair) error {
	_, err := service.UserStorage.GetByGmail(user.Gmail)

	if err != nil {
		return err
	}

	return service.TemporaryStorage.Create(user)
}

func (service *ResetPasswordService) ChangeUserPassword(user entities.GmailWithKeyPair, newPassword string) error {
	_, err := service.TemporaryStorage.GetByUniqueKey(user.Key)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	err = service.TemporaryStorage.Delete(user)
	if err != nil {
		return fmt.Errorf("could not remove user")
	}

	realUser, err := service.UserStorage.GetByGmail(user.Gmail)
	if err != nil {
		return err
	}

	return service.UserStorage.UpdatePassword(realUser, newPassword)
}
