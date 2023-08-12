package signin

import (
	"auth/src/domain/entities"
	"auth/src/domain/errors"
	"auth/src/utils"
)

type repository interface {
	ByID(string) (entities.User, error)
	ReadAll() ([]entities.User, error)
}

type signinService struct {
	db   repository
	hash func(string) string
}

func NewSigninService(db repository, hash func(string) string) signinService {
	return signinService{db: db, hash: hash}
}

func (s signinService) Login(gmail, password string) (entities.User, error) {
	users, err := s.db.ReadAll()
	if err != nil {
		return entities.User{}, err
	}

	user, err := utils.Find(users, func(u entities.User) bool {
		return u.Gmail == gmail && u.Password == s.hash(password)
	})
	if err != nil {
		return entities.User{}, errors.ErrUserNotFound
	}

	return user, nil
}
