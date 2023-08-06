package signin

import (
	"auth/src/domain/entities"
	"auth/src/domain/errors"
	"auth/src/utils"
)

type repository interface {
	ByID(string) (entities.UserEntity, error)
	ReadAll() ([]entities.UserEntity, error)
}

type signinService struct {
	db   repository
	hash func(string) string
}

func NewSigninService(db repository, hash func(string) string) signinService {
	return signinService{db: db, hash: hash}
}

func (s signinService) Login(gmail, password string) (entities.UserEntity, error) {
	users, err := s.db.ReadAll()
	if err != nil {
		return entities.UserEntity{}, err
	}

	user, err := utils.Find(users, func(u entities.UserEntity) bool {
		return u.Gmail == gmail && u.Password == s.hash(password)
	})
	if err != nil {
		return entities.UserEntity{}, errors.ErrUserNotFound
	}

	return user, nil
}
