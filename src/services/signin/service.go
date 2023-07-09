package signin

import (
	"auth/src/entities"
	"auth/src/errors"
	"auth/src/storage"
)

type repository interface {
	ReadAll() ([]entities.UserEntity, error)
}

type signinService struct {
	db   repository
	hash func(string) string
}

func NewSigninService(db repository, hash func(string) string) signinService {
	return signinService{db: db, hash: hash}
}

func (s signinService) Login(gmail, password string) (string, error) {
	users, err := s.db.ReadAll()
	if err != nil {
		return "", err
	}

	user, err := storage.Find(users, func(u entities.UserEntity) bool {
		return u.Gmail == gmail && u.Password == s.hash(password)
	})
	if err != nil {
		return "", errors.ErrUserNotFound
	}

	return user.ID, nil
}
