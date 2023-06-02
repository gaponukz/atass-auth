package storage

import (
	"auth/src/entities"
	"fmt"
)

type MemoryStorage struct{}

func (stor *MemoryStorage) Create(entities.User) error {
	return nil
}

func (stor *MemoryStorage) Update(entities.User) error {
	return nil
}

func (stor *MemoryStorage) ReadAll() ([]entities.User, error) {
	users := []entities.User{
		{Gmail: "user1@gmail.com", Password: "12345", FullName: "Anna Dou"},
		{Gmail: "user2@gmail.com", Password: "hvdavf", FullName: "Vlad Feq"},
		{Gmail: "user3@gmail.com", Password: "ahevre", FullName: "Alex Ogh"},
		{Gmail: "user4@gmail.com", Password: "3qbduag3", FullName: "Max Daz"},
	}

	return users, nil
}

func (stor *MemoryStorage) Delete(entities.User) error {
	return nil
}

func (stor *MemoryStorage) GetByGmail(gmail string) (entities.User, error) {
	var userId int = -1
	users, err := stor.ReadAll()

	if err != nil {
		return entities.User{}, err
	}

	for idx, user := range users {
		if user.Gmail == gmail {
			userId = idx
			break
		}
	}

	if userId == -1 {
		return entities.User{}, fmt.Errorf("user %s not found", gmail)
	}

	return users[userId], nil
}
