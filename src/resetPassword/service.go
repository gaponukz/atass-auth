package resetPassword

import (
	"auth/src/entities"
	"auth/src/storage"
	"fmt"
	"math/rand"
	"strconv"
)

type UpdatePasswordGetByGmailAbleStorage interface {
	GetByGmail(string) (entities.User, error)
	UpdatePassword(entities.User, string) error
}

type ResetPasswordService struct {
	TemporaryStorage storage.ITemporaryStorage[storage.UserCredentials]
	UserStorage      UpdatePasswordGetByGmailAbleStorage
	Notify           func(gmail, key string) error
}

func (service *ResetPasswordService) GenerateAndSendCodeToGmail(userGmail string) (string, error) {
	key := strconv.Itoa(rand.Intn(900000) + 100000)
	err := service.Notify(userGmail, key)

	return key, err
}

func (service *ResetPasswordService) AddUserToTemporaryStorage(user storage.UserCredentials) error {
	_, err := service.UserStorage.GetByGmail(user.Gmail)

	if err != nil {
		return err
	}

	return service.TemporaryStorage.Create(user)
}

func (service *ResetPasswordService) ChangeUserPassword(user storage.UserCredentials, newPassword string) error {
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
