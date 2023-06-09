package registration

import (
	"auth/src/entities"
	"auth/src/storage"
	"fmt"
	"math/rand"
	"strconv"
)

type CreateAndGetByGmailAbleStorage interface {
	Create(entities.User) error
	GetByGmail(string) (entities.User, error)
}

type RegistrationService struct {
	UserStorage       CreateAndGetByGmailAbleStorage
	FutureUserStorage storage.ITemporaryStorage[entities.FutureUser]
	Notify            func(gmail, key string) error
}

func (service *RegistrationService) generateKey() string {
	return strconv.Itoa(rand.Intn(900000) + 100000)
}

func (service *RegistrationService) GetInformatedFutureUser(userGmail string) (string, error) {
	key := service.generateKey()
	err := service.Notify(userGmail, key)

	return key, err
}

func (service *RegistrationService) AddUserToFutureStorage(user entities.FutureUser) error {
	mayUser, _ := service.UserStorage.GetByGmail(user.Gmail)

	if mayUser.Gmail != "" {
		return fmt.Errorf("already registered gmail")
	}

	return service.FutureUserStorage.Create(user)
}

func (service *RegistrationService) RemoveUserFromFutureStorage(user entities.FutureUser) error {
	_, err := service.FutureUserStorage.GetByUniqueKey(user.UniqueKey)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	err = service.FutureUserStorage.Delete(user)
	if err != nil {
		return fmt.Errorf("could not remove user")
	}

	return service.UserStorage.Create(entities.User{
		Gmail:    user.Gmail,
		Password: user.Password,
		FullName: user.FullName,
		Phone:    user.Phone,
	})
}
