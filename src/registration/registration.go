package registration

import (
	"auth/src/entities"
	"auth/src/storage"
	"fmt"
)

type createAndReadAbleStorage interface {
	Create(entities.User) (entities.UserEntity, error)
	ReadAll() ([]entities.UserEntity, error)
}

type gmailKeyPairStorage interface {
	Create(entities.GmailWithKeyPair) error
	Delete(entities.GmailWithKeyPair) error
	GetByUniqueKey(string) (entities.GmailWithKeyPair, error)
}

func NewRegistrationService(
	userStorage createAndReadAbleStorage,
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
	userStorage       createAndReadAbleStorage
	futureUserStorage gmailKeyPairStorage
	notify            func(gmail, key string) error
	generateCode      func() string
}

func (s registrationService) SendGeneratedCode(userGmail string) (string, error) {
	key := s.generateCode()
	err := s.notify(userGmail, key) // TODO: make gorutine with 5 sec deadline context

	return key, err
}

func (s registrationService) AddUserToTemporaryStorage(user entities.GmailWithKeyPair) error {
	users, err := s.userStorage.ReadAll()
	if err != nil {
		return err
	}

	isExist := storage.IsExist(users, func(u entities.UserEntity) bool {
		return u.Gmail == user.Gmail
	})

	if isExist {
		return fmt.Errorf("already registered gmail")
	}

	return s.futureUserStorage.Create(user)
}

func (s registrationService) RegisterUserOnRightCode(pair entities.GmailWithKeyPair, user entities.User) (string, error) {
	_, err := s.futureUserStorage.GetByUniqueKey(pair.Key)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	err = s.futureUserStorage.Delete(pair)
	if err != nil {
		return "", fmt.Errorf("could not remove user")
	}

	newUser, err := s.userStorage.Create(user)
	if err != nil {
		return "", err
	}

	return newUser.ID, nil
}
