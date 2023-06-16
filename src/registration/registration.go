package registration

import (
	"auth/src/entities"
	"fmt"
)

type createAndGetByGmailAbleStorage interface {
	Create(entities.User) error
	GetByGmail(string) (entities.User, error)
}

type gmailKeyPairStorage interface {
	Create(entities.GmailWithKeyPair) error
	Delete(entities.GmailWithKeyPair) error
	GetByUniqueKey(string) (entities.GmailWithKeyPair, error)
}

type RegistrationService struct {
	UserStorage       createAndGetByGmailAbleStorage
	FutureUserStorage gmailKeyPairStorage
	Notify            func(gmail, key string) error
	GenerateCode      func() string
}

func (service *RegistrationService) GetInformatedFutureUser(userGmail string) (string, error) {
	key := service.GenerateCode()
	err := service.Notify(userGmail, key) // TODO: make gorutine with 5 sec deadline context

	return key, err
}

func (service *RegistrationService) AddUserToTemporaryStorage(user entities.GmailWithKeyPair) error {
	mayUser, _ := service.UserStorage.GetByGmail(user.Gmail)

	if mayUser.Gmail != "" {
		return fmt.Errorf("already registered gmail")
	}

	return service.FutureUserStorage.Create(user)
}

func (service *RegistrationService) RegisterUserOnRightCode(pair entities.GmailWithKeyPair, user entities.User) error {
	_, err := service.FutureUserStorage.GetByUniqueKey(pair.Key)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	err = service.FutureUserStorage.Delete(pair)
	if err != nil {
		return fmt.Errorf("could not remove user")
	}

	return service.UserStorage.Create(user)
}
