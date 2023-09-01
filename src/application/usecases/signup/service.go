package signup

import (
	"auth/src/application/dto"
	"auth/src/domain/entities"
	"auth/src/domain/errors"
	"auth/src/utils"
	"regexp"
	"strings"
	"unicode"
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

func (s registrationService) IsPasswordValid(password string) bool {
	const minLength = 8
	var hasUppercase, hasLowercase, hasSpecialChar bool

	specialCharRegex := regexp.MustCompile(`[!@#$%^&*()-_+=\[\]{}|:;"'<>,.?/~]`)

	if len(password) < minLength {
		return false
	}

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
		}
		if unicode.IsLower(char) {
			hasLowercase = true
		}
		if specialCharRegex.MatchString(string(char)) || unicode.IsDigit(char) {
			hasSpecialChar = true
		}
	}

	return hasUppercase && hasLowercase && hasSpecialChar
}

func (s registrationService) IsPhoneNumberValid(phoneNumber string) bool {
	phoneRegex := regexp.MustCompile(`^(\+)?(\d+\s?)+$`)

	if len(strings.Replace(phoneNumber, " ", "", -1)) > 15 {
		return false
	}

	return phoneRegex.MatchString(phoneNumber)
}

func (s registrationService) IsFullNameValid(fullName string) bool {
	fullNameRegex := regexp.MustCompile(`^[a-zA-Z]{2,} [a-zA-Z]{2,}$`)
	return fullNameRegex.MatchString(fullName)
}
