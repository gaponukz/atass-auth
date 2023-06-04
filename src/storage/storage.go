package storage

import "auth/src/entities"

type IUserStorage interface {
	Create(entities.User) error
	ReadAll() ([]entities.User, error)
	Delete(entities.User) error
	GetByGmail(string) (entities.User, error)
}

type IFutureUserStorage interface {
	Create(entities.FutureUser) error
	Delete(entities.FutureUser) error
	GetByUniqueKey(string) (entities.FutureUser, error)
}
