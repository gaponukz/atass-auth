package utils

import (
	"auth/src/entities"
	"errors"
)

func Filter(users []entities.UserEntity, filterFunc func(user entities.UserEntity) bool) []entities.UserEntity {
	var filteredUsers []entities.UserEntity
	for _, user := range users {
		if filterFunc(user) {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers
}

func Find(users []entities.UserEntity, filterFunc func(user entities.UserEntity) bool) (entities.UserEntity, error) {
	for _, user := range users {
		if filterFunc(user) {
			return user, nil
		}
	}

	return entities.UserEntity{}, errors.New("not found")
}

func IsExist(users []entities.UserEntity, filterFunc func(user entities.UserEntity) bool) bool {
	for _, user := range users {
		if filterFunc(user) {
			return true
		}
	}

	return false
}
