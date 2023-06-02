package storage

import "auth/src/entities"

type IUserStorage interface {
	Create(entities.User) error
	Update(entities.User) error
	ReadAll() ([]entities.User, error)
	Delete(entities.User) error

	GetByGmail(string) (entities.User, error)
}
