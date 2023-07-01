package mocks

import (
	"auth/src/entities"
	"errors"
	"fmt"
)

type temporaryStorageMock struct {
	users []entities.GmailWithKeyPair
}

func NewTemporaryStorageMock() *temporaryStorageMock {
	return &temporaryStorageMock{}
}

func (m *temporaryStorageMock) Create(pair entities.GmailWithKeyPair) error {
	m.users = append(m.users, pair)
	return nil
}

func (m *temporaryStorageMock) Delete(pair entities.GmailWithKeyPair) error {
	for i, user := range m.users {
		if user.Gmail == pair.Gmail && user.Key == pair.Key {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return errors.New("pair not found")
}

func (m *temporaryStorageMock) GetByUniqueKey(key string) (entities.GmailWithKeyPair, error) {
	for _, user := range m.users {
		if user.Key == key {
			return user, nil
		}
	}
	return entities.GmailWithKeyPair{}, errors.New("pair not found")
}

type mockStorage struct {
	users []entities.UserEntity
}

func (m *mockStorage) Create(user entities.User) (entities.UserEntity, error) {
	newUser := entities.UserEntity{ID: "1", User: user}
	m.users = append(m.users, newUser)
	return newUser, nil
}

func (m *mockStorage) ReadAll() ([]entities.UserEntity, error) {
	return m.users, nil
}

func (m *mockStorage) Update(userToUpdate entities.UserEntity) error {
	for idx, user := range m.users {
		if user.ID == userToUpdate.ID {
			m.users[idx] = userToUpdate
			return nil
		}
	}

	return fmt.Errorf("user %s not found", userToUpdate.ID)
}

func NewMockStorage() *mockStorage {
	return &mockStorage{}
}
