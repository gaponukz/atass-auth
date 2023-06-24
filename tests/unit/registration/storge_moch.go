package registration

import (
	"auth/src/entities"
	"errors"
)

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

func NewMockStorage() *mockStorage {
	return &mockStorage{}
}
