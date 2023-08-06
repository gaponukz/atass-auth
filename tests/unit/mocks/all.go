package mocks

import (
	"auth/src/application/dto"
	"auth/src/domain/entities"
	"auth/src/domain/errors"
	"fmt"
)

type temporaryStorageMock struct {
	users []dto.GmailWithKeyPairDTO
}

func NewTemporaryStorageMock() *temporaryStorageMock {
	return &temporaryStorageMock{}
}

func (m *temporaryStorageMock) Create(pair dto.GmailWithKeyPairDTO) error {
	m.users = append(m.users, pair)
	return nil
}

func (m *temporaryStorageMock) Delete(pair dto.GmailWithKeyPairDTO) error {
	for i, user := range m.users {
		if user.Gmail == pair.Gmail && user.Key == pair.Key {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("404")
}

func (m *temporaryStorageMock) GetByUniqueKey(key string) (dto.GmailWithKeyPairDTO, error) {
	for _, user := range m.users {
		if user.Key == key {
			return user, nil
		}
	}
	return dto.GmailWithKeyPairDTO{}, fmt.Errorf("404")
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

	return errors.ErrUserNotFound
}

func (m *mockStorage) ByID(id string) (entities.UserEntity, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}

	return entities.UserEntity{}, errors.ErrUserNotFound
}

func NewMockStorage() *mockStorage {
	return &mockStorage{}
}
