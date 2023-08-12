package utils

import (
	"auth/src/domain/entities"
	"errors"
)

func Filter(users []entities.User, filterFunc func(user entities.User) bool) []entities.User {
	var filteredUsers []entities.User
	for _, user := range users {
		if filterFunc(user) {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers
}

func Find(users []entities.User, filterFunc func(user entities.User) bool) (entities.User, error) {
	for _, user := range users {
		if filterFunc(user) {
			return user, nil
		}
	}

	return entities.User{}, errors.New("not found")
}

func IsExist(users []entities.User, filterFunc func(user entities.User) bool) bool {
	for _, user := range users {
		if filterFunc(user) {
			return true
		}
	}

	return false
}
