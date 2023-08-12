package signup

import (
	"auth/src/application/dto"
	"auth/src/domain/entities"
	"auth/src/domain/errors"
	"auth/src/utils"
)

type notifier interface {
	Notify(string, string) error
}

type createAndReadAbleStorage interface {
	Create(dto.CreateUserDTO) (entities.User, error)
	ReadAll() ([]entities.User, error)
}

type gmailKeyPairStorage interface {
	Create(dto.GmailWithKeyPairDTO) error
	Delete(dto.GmailWithKeyPairDTO) error
	GetByUniqueKey(string) (dto.GmailWithKeyPairDTO, error)
}

func NewRegistrationService(
	userStorage createAndReadAbleStorage,
	futureUserStorage gmailKeyPairStorage,
	notifier notifier,
	generateCode func() string,
	hash func(string) string,
) *registrationService {
	return &registrationService{
		userStorage:       userStorage,
		futureUserStorage: futureUserStorage,
		notifier:          notifier,
		generateCode:      generateCode,
		hash:              hash,
	}
}

type registrationService struct {
	userStorage       createAndReadAbleStorage
	futureUserStorage gmailKeyPairStorage
	notifier          notifier
	generateCode      func() string
	hash              func(string) string
}

func (s registrationService) SendGeneratedCode(userGmail string) (string, error) {
	users, err := s.userStorage.ReadAll()
	if err != nil {
		return "", err
	}

	isExist := utils.IsExist(users, func(u entities.User) bool {
		return u.Gmail == userGmail
	})

	if isExist {
		return "", errors.ErrUserAlreadyExists
	}

	key := s.generateCode()
	err = s.notifier.Notify(userGmail, key)

	return key, err
}

func (s registrationService) AddUserToTemporaryStorage(user dto.GmailWithKeyPairDTO) error {
	return s.futureUserStorage.Create(user)
}

func (s registrationService) RegisterUserOnRightCode(user dto.SignUpDTO) (string, error) {
	err := s.futureUserStorage.Delete(dto.GmailWithKeyPairDTO{Gmail: user.Gmail, Key: user.Key})
	if err != nil {
		return "", errors.ErrRegisterRequestMissing
	}

	user.Password = s.hash(user.Password)
	newUser, err := s.userStorage.Create(dto.CreateUserDTO{
		Gmail:               user.Gmail,
		Phone:               user.Phone,
		FullName:            user.FullName,
		Password:            user.Password,
		AllowsAdvertisement: user.AllowsAdvertisement,
	})
	if err != nil {
		return "", err
	}

	return newUser.ID, nil
}
