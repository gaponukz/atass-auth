package mocks

import (
	"auth/src/entities"
	"errors"
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
	users []entities.User
}

func (m *mockStorage) Create(user entities.User) error {
	m.users = append(m.users, user)
	return nil
}

func (m *mockStorage) GetByGmail(gmail string) (entities.User, error) {
	for _, user := range m.users {
		if user.Gmail == gmail {
			return user, nil
		}
	}

	return entities.User{}, errors.New("User not found")
}

func (m *mockStorage) UpdatePassword(user entities.User, password string) error {
	for i, u := range m.users {
		if u.Gmail == user.Gmail {
			m.users = append(m.users[:i], m.users[i+1:]...)

			updatedUser := user
			updatedUser.Password = password

			m.users = append(m.users, updatedUser)
			return nil
		}
	}
	return errors.New("User not found")
}

func NewMockStorage() *mockStorage {
	return &mockStorage{}
}
