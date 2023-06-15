package resetPassword

import (
	"auth/src/entities"
	"fmt"
	"math/rand"
	"strconv"
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
}

func (service *ResetPasswordService) GenerateAndSendCodeToGmail(userGmail string) (string, error) {
	key := strconv.Itoa(rand.Intn(900000) + 100000)
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
