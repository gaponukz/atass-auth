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

func NewRegistrationService(
	userStorage createAndGetByGmailAbleStorage,
	futureUserStorage gmailKeyPairStorage,
	notify func(gmail, key string) error,
	generateCode func() string,
) *registrationService {
	return &registrationService{
		userStorage:       userStorage,
		futureUserStorage: futureUserStorage,
		notify:            notify,
		generateCode:      generateCode,
	}
}

type registrationService struct {
	userStorage       createAndGetByGmailAbleStorage
	futureUserStorage gmailKeyPairStorage
	notify            func(gmail, key string) error
	generateCode      func() string
}

func (s registrationService) GetInformatedFutureUser(userGmail string) (string, error) {
	key := s.generateCode()
	err := s.notify(userGmail, key) // TODO: make gorutine with 5 sec deadline context

	return key, err
}

func (s registrationService) AddUserToTemporaryStorage(user entities.GmailWithKeyPair) error {
	mayUser, _ := s.userStorage.GetByGmail(user.Gmail)

	if mayUser.Gmail != "" {
		return fmt.Errorf("already registered gmail")
	}

	return s.futureUserStorage.Create(user)
}

func (s registrationService) RegisterUserOnRightCode(pair entities.GmailWithKeyPair, user entities.User) error {
	_, err := s.futureUserStorage.GetByUniqueKey(pair.Key)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	err = s.futureUserStorage.Delete(pair)
	if err != nil {
		return fmt.Errorf("could not remove user")
	}

	return s.userStorage.Create(user)
}
