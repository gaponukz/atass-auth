package registration

import (
	"auth/src/entities"
	"auth/src/errors"
	"auth/src/storage"
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
		return errors.ErrUserAlreadyExists
	}

	return s.futureUserStorage.Create(user)
}

func (s registrationService) RegisterUserOnRightCode(pair entities.GmailWithKeyPair, user entities.User) (string, error) {
	err := s.futureUserStorage.Delete(pair)
	if err != nil {
		return "", errors.ErrRegisterRequestMissing
	}

	newUser, err := s.userStorage.Create(user)
	if err != nil {
		return "", err
	}

	return newUser.ID, nil
}
