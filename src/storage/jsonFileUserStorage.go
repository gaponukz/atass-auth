package storage

import (
	"auth/src/entities"
	"auth/src/security"
	"encoding/json"
	"fmt"
	"os"
)

type userJsonFileStorage struct {
	filePath string
}

func NewUserJsonFileStorage(filePath string) *userJsonFileStorage {
	return &userJsonFileStorage{filePath: filePath}
}

func (s userJsonFileStorage) Create(user entities.User) error {
	users, err := s.readUsersFromFile()
	if err != nil {
		return err
	}

	user.Password = security.GetSha256(user.Password)

	users = append(users, user)
	err = s.writeUsersToFile(users)
	if err != nil {
		return err
	}

	return nil
}

func (s userJsonFileStorage) Delete(userToRemove entities.User) error {
	users, err := s.readUsersFromFile()
	if err != nil {
		return err
	}

	index := -1

	for idx, user := range users {
		if user.Gmail == userToRemove.Gmail {
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

func (s userJsonFileStorage) GetByGmail(gmail string) (entities.User, error) {
	users, err := s.readUsersFromFile()
	if err != nil {
		return entities.User{}, err
	}

	var userFound entities.User
	userFoundIndex := -1

	for idx, user := range users {
		if user.Gmail == gmail {
			userFound = user
			userFoundIndex = idx
			break
		}
	}

	if userFoundIndex == -1 {
		return entities.User{}, fmt.Errorf("user %s not found", gmail)
	}

	return userFound, nil
}

func (s userJsonFileStorage) readUsersFromFile() ([]entities.User, error) {
	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []entities.User
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s userJsonFileStorage) UpdatePassword(userToUpdate entities.User, newPassword string) error {
	users, err := s.readUsersFromFile()
	if err != nil {
		return err
	}

	for idx, user := range users {
		if user.Gmail == userToUpdate.Gmail {
			users[idx].Password = security.GetSha256(newPassword)
			err = s.writeUsersToFile(users)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("user %s not found", userToUpdate.Gmail)
}

func (s userJsonFileStorage) AddSubscribedRoute(userToUpdate entities.User, routeId string) error {
	users, err := s.readUsersFromFile()
	if err != nil {
		return err
	}

	for idx, user := range users {
		if user.Gmail == userToUpdate.Gmail {
			users[idx].PurchasedRouteIds = append(users[idx].PurchasedRouteIds, routeId)
			err = s.writeUsersToFile(users)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("user %s not found", userToUpdate.Gmail)
}

func (s userJsonFileStorage) writeUsersToFile(users []entities.User) error {
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
