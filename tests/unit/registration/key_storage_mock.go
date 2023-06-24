package registration

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
