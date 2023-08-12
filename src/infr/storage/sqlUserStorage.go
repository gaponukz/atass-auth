package storage

import (
	"auth/src/application/dto"
	"auth/src/domain/entities"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbProxyUser struct {
	ID                  string `gorm:"primaryKey;not null;unique"`
	Gmail               string `gorm:"not null;unique"`
	Password            string `gorm:"not null;unique"`
	Phone               string `gorm:"not null;unique"`
	FullName            string
	AllowsAdvertisement bool
	PurchasedRouteIds   string
}

type sqlUserStorage struct {
	db *gorm.DB
}

func fromProxyToUser(p dbProxyUser) entities.User {
	var pathes []entities.Path

	err := json.Unmarshal([]byte(p.PurchasedRouteIds), &pathes)
	if err != nil {
		pathes = nil
	}

	return entities.User{
		ID:                  p.ID,
		Gmail:               p.Gmail,
		Password:            p.Password,
		Phone:               p.Phone,
		FullName:            p.FullName,
		AllowsAdvertisement: p.AllowsAdvertisement,
		PurchasedRouteIds:   pathes,
	}
}

func fromUserToProxy(u entities.User) dbProxyUser {
	var pathes string = "[]"

	bPathes, err := json.Marshal(u.PurchasedRouteIds)
	if err == nil {
		pathes = string(bPathes)
	}

	return dbProxyUser{
		ID:                  u.ID,
		Gmail:               u.Gmail,
		Password:            u.Password,
		Phone:               u.Phone,
		FullName:            u.FullName,
		AllowsAdvertisement: u.AllowsAdvertisement,
		PurchasedRouteIds:   pathes,
	}
}

type PostgresCredentials struct {
	Host     string
	User     string
	Password string
	Dbname   string
	Port     string
	Sslmode  string
}

func NewPostgresUserStorage(c PostgresCredentials) (*sqlUserStorage, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host, c.User, c.Password, c.Dbname, c.Port, c.Sslmode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&dbProxyUser{})
	if err != nil {
		return nil, err
	}

	return &sqlUserStorage{
		db: db,
	}, nil
}

func (repo sqlUserStorage) Create(createDto dto.CreateUserDTO) (entities.User, error) {
	entity := entities.User{
		ID:                  uuid.New().String(),
		Gmail:               createDto.Gmail,
		Password:            createDto.Password,
		Phone:               createDto.Phone,
		FullName:            createDto.FullName,
		AllowsAdvertisement: createDto.AllowsAdvertisement,
	}
	proxy := fromUserToProxy(entity)

	result := repo.db.Create(&proxy)
	if result.Error != nil {
		return entities.User{}, fmt.Errorf("can not create user %s: %v", createDto.Gmail, result.Error)
	}

	return entity, nil
}

func (repo sqlUserStorage) ReadAll() ([]entities.User, error) {
	var proxies []dbProxyUser

	result := repo.db.Find(&proxies)
	if result.Error != nil {
		return nil, fmt.Errorf("can not get all users: %v", result.Error)
	}

	users := make([]entities.User, len(proxies))

	for i, proxy := range proxies {
		users[i] = fromProxyToUser(proxy)
	}

	return users, nil
}

func (repo sqlUserStorage) ByID(id string) (entities.User, error) {
	var proxy dbProxyUser

	result := repo.db.Where("id = ?", id).First(&proxy)
	if result.Error != nil {
		return entities.User{}, fmt.Errorf("can not get by id %s: %v", id, result.Error)
	}

	return fromProxyToUser(proxy), nil
}

func (repo sqlUserStorage) Update(userToUpdate entities.User) error {
	proxy := fromUserToProxy(userToUpdate)

	result := repo.db.Save(&proxy)
	if result.Error != nil {
		return fmt.Errorf("can not update %s: %v", userToUpdate.ID, result.Error)
	}

	return nil
}

func (repo sqlUserStorage) Delete(id string) error {
	result := repo.db.Where("id = ?", id).Delete(&dbProxyUser{})

	if result.Error != nil {
		return fmt.Errorf("can not delete %s: %v", id, result.Error)
	}

	return nil
}

func (repo sqlUserStorage) DropTable() error {
	return repo.db.Migrator().DropTable(&dbProxyUser{})
}
