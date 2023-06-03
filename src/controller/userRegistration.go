package controller

import (
	"auth/src/entities"
	"fmt"
)

type CreateAndGetByGmailStorage interface {
	Create(entities.User) error
	GetByGmail(string) (entities.User, error)
}

func registerUser(userDTO userCredentialsnDTO, storage CreateAndGetByGmailStorage) error {
	user, _ := storage.GetByGmail(userDTO.Gmail)

	if user.Gmail != "" {
		return fmt.Errorf("already registered gmail")
	}

	return storage.Create(entities.User{
		Gmail:    userDTO.Gmail,
		Password: userDTO.Password,
		FullName: userDTO.FullName,
	})
}
