package storage

import (
	"auth/src/domain/entities"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type userJsonFileStorage struct {
	filePath string
}

func NewUserJsonFileStorage(filePath string) *userJsonFileStorage {
	return &userJsonFileStorage{filePath: filePath}
}

func (s userJsonFileStorage) Create(user entities.User) (entities.UserEntity, error) {
	users, err := s.readUsersFromFile()
	if err != nil {
		return entities.UserEntity{}, err
	}

	userEntity := entities.UserEntity{ID: uuid.New().String(), User: user}
	users = append(users, userEntity)
	err = s.writeUsersToFile(users)

	if err != nil {
		return entities.UserEntity{}, err
	}

	return userEntity, nil
}

func (s userJsonFileStorage) ReadAll() ([]entities.UserEntity, error) {
	return s.readUsersFromFile()
}

func (s userJsonFileStorage) ByID(id string) (entities.UserEntity, error) {
	users, err := s.readUsersFromFile()
	if err != nil {
		return entities.UserEntity{}, err
	}

	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return entities.UserEntity{}, fmt.Errorf("user %s not found", id)
}

func (s userJsonFileStorage) Update(userToUpdate entities.UserEntity) error {
	users, err := s.readUsersFromFile()
	if err != nil {
		return err
	}

	for idx, user := range users {
		if user.ID == userToUpdate.ID {
			users[idx] = userToUpdate
			err = s.writeUsersToFile(users)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("user %s not found", userToUpdate.ID)
}

func (s userJsonFileStorage) Delete(id string) error {
	users, err := s.readUsersFromFile()
	if err != nil {
		return err
	}

	index := -1

	for idx, user := range users {
		if user.ID == id {
			index = idx
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("not found")
	}

	users = append(users[:index], users[index+1:]...)
	err = s.writeUsersToFile(users)
	if err != nil {
		return err
	}

	return nil
}

func (s userJsonFileStorage) readUsersFromFile() ([]entities.UserEntity, error) {
	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []entities.UserEntity
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s userJsonFileStorage) writeUsersToFile(users []entities.UserEntity) error {
	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(users)
	if err != nil {
		return err
	}

	return nil
}
