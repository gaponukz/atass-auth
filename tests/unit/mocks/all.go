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
	users []entities.User
}

func (m *mockStorage) Create(createDto dto.CreateUserDTO) (entities.User, error) {
	newUser := entities.User{
		ID:                  "1",
		Gmail:               createDto.Gmail,
		Password:            createDto.Password,
		Phone:               createDto.Phone,
		FullName:            createDto.FullName,
		AllowsAdvertisement: createDto.AllowsAdvertisement,
	}
	m.users = append(m.users, newUser)
	return newUser, nil
}

func (m *mockStorage) ReadAll() ([]entities.User, error) {
	return m.users, nil
}

func (m *mockStorage) Update(userToUpdate entities.User) error {
	for idx, user := range m.users {
		if user.ID == userToUpdate.ID {
			m.users[idx] = userToUpdate
			return nil
		}
	}

	return errors.ErrUserNotFound
}

func (m *mockStorage) ByID(id string) (entities.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}

	return entities.User{}, errors.ErrUserNotFound
}

func NewMockStorage() *mockStorage {
	return &mockStorage{}
}

type mockUserNotifier struct{}

func (m mockUserNotifier) NotifyUser(user entities.User, mes string) error {
	return nil
}

func NewMockUserNotifier() mockUserNotifier {
	return mockUserNotifier{}
}

type mockGmailNotifier struct{}

func (m mockGmailNotifier) Notify(user string, mes string) error {
	return nil
}

func NewMockGmailNotifier() mockGmailNotifier {
	return mockGmailNotifier{}
}
