package storage

import (
	"auth/src/entities"
	"fmt"
)

type UserMemoryStorage struct {
	Users []entities.User
}

func (strg *UserMemoryStorage) Create(user entities.User) error {
	strg.Users = append(strg.Users, user)
	return nil
}

func (stor *UserMemoryStorage) ReadAll() ([]entities.User, error) {
	return stor.Users, nil
}

func (strg *UserMemoryStorage) Delete(userToRemove entities.User) error {
	index := -1

	for idx, user := range strg.Users {
		if user.Gmail == userToRemove.Gmail {
			index = idx
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("not found")
	}

	strg.Users = append(strg.Users[:index], strg.Users[index+1])

	return nil
}

func (stor *UserMemoryStorage) GetByGmail(gmail string) (entities.User, error) {
	var userId int = -1
	users, err := stor.ReadAll()

	if err != nil {
		return entities.User{}, err
	}

	for idx, user := range users {
		if user.Gmail == gmail {
			userId = idx
			break
		}
	}

	if userId == -1 {
		return entities.User{}, fmt.Errorf("user %s not found", gmail)
	}

	return users[userId], nil
}
