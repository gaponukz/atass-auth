package signup

import (
	"auth/src/dto"
	"auth/src/entities"
	"auth/src/errors"
	"auth/src/utils"
)

type createAndReadAbleStorage interface {
	Create(entities.User) (entities.UserEntity, error)
	ReadAll() ([]entities.UserEntity, error)
}

type gmailKeyPairStorage interface {
	Create(dto.GmailWithKeyPairDTO) error
	Delete(dto.GmailWithKeyPairDTO) error
	GetByUniqueKey(string) (dto.GmailWithKeyPairDTO, error)
}

func NewRegistrationService(
	userStorage createAndReadAbleStorage,
	futureUserStorage gmailKeyPairStorage,
	notify func(gmail, key string) error,
	generateCode func() string,
	hash func(string) string,
) *registrationService {
	return &registrationService{
		userStorage:       userStorage,
		futureUserStorage: futureUserStorage,
		notify:            notify,
		generateCode:      generateCode,
		hash:              hash,
	}
}

type registrationService struct {
	userStorage       createAndReadAbleStorage
	futureUserStorage gmailKeyPairStorage
	notify            func(gmail, key string) error
	generateCode      func() string
	hash              func(string) string
}

func (s registrationService) SendGeneratedCode(userGmail string) (string, error) {
	key := s.generateCode()
	err := s.notify(userGmail, key) // TODO: make gorutine with 5 sec deadline context

	return key, err
}

func (s registrationService) AddUserToTemporaryStorage(user dto.GmailWithKeyPairDTO) error {
	users, err := s.userStorage.ReadAll()
	if err != nil {
		return err
	}

	isExist := utils.IsExist(users, func(u entities.UserEntity) bool {
		return u.Gmail == user.Gmail
	})

	if isExist {
		return errors.ErrUserAlreadyExists
	}

	return s.futureUserStorage.Create(user)
}

func (s registrationService) RegisterUserOnRightCode(user dto.SignUpDTO) (string, error) {
	err := s.futureUserStorage.Delete(dto.GmailWithKeyPairDTO{Gmail: user.Gmail, Key: user.Key})
	if err != nil {
		return "", errors.ErrRegisterRequestMissing
	}

	user.Password = s.hash(user.Password)
	newUser, err := s.userStorage.Create(entities.User{
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
