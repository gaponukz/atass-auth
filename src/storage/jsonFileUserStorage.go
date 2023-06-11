package storage

import (
	"auth/src/entities"
	"encoding/json"
	"fmt"
	"os"
)

type UserJsonFileStorage struct {
	FilePath string
}

func (strg *UserJsonFileStorage) Create(user entities.User) error {
	users, err := strg.readUsersFromFile()
	if err != nil {
		return err
	}

	user.Password = GetSha256(user.Password)

	users = append(users, user)
	err = strg.writeUsersToFile(users)
	if err != nil {
		return err
	}

	return nil
}

func (strg *UserJsonFileStorage) Delete(userToRemove entities.User) error {
	users, err := strg.readUsersFromFile()
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
	err = strg.writeUsersToFile(users)
	if err != nil {
		return err
	}

	return nil
}

func (stor *UserJsonFileStorage) GetByGmail(gmail string) (entities.User, error) {
	users, err := stor.readUsersFromFile()
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

func (strg *UserJsonFileStorage) readUsersFromFile() ([]entities.User, error) {
	file, err := os.Open(strg.FilePath)
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

func (strg *UserJsonFileStorage) UpdatePassword(userToUpdate entities.User, newPassword string) error {
	users, err := strg.readUsersFromFile()
	if err != nil {
		return err
	}

	for idx, user := range users {
		if user.Gmail == userToUpdate.Gmail {
			users[idx].Password = GetSha256(newPassword)
			err = strg.writeUsersToFile(users)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("user %s not found", userToUpdate.Gmail)
}

func (strg *UserJsonFileStorage) SubscribeToTheRoute(userToUpdate entities.User, routeId string) error {
	users, err := strg.readUsersFromFile()
	if err != nil {
		return err
	}

	for idx, user := range users {
		if user.Gmail == userToUpdate.Gmail {
			users[idx].PurchasedRouteIds = append(users[idx].PurchasedRouteIds, routeId)
			err = strg.writeUsersToFile(users)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("user %s not found", userToUpdate.Gmail)
}

func (strg *UserJsonFileStorage) writeUsersToFile(users []entities.User) error {
	file, err := os.Create(strg.FilePath)
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
