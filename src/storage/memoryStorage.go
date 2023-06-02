package storage

import "auth/src/entities"

type MemoryStorage struct{}

func (stor *MemoryStorage) Create(entities.User) error {
	return nil
}

func (stor *MemoryStorage) Update(entities.User) error {
	return nil
}

func (stor *MemoryStorage) ReadAll() ([]entities.User, error) {
	users := []entities.User{
		entities.User{Gmail: "user1@gmail.com", Password: "12345", FullName: "Anna Dou"},
		entities.User{Gmail: "user2@gmail.com", Password: "hvdavf", FullName: "Vlad Feq"},
		entities.User{Gmail: "user3@gmail.com", Password: "ahevre", FullName: "Alex Ogh"},
		entities.User{Gmail: "user4@gmail.com", Password: "3qbduag3", FullName: "Max Daz"},
	}

	return users, nil
}

func (stor *MemoryStorage) Delete(entities.User) error {
	return nil
}
