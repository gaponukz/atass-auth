package storage

import (
	"auth/src/entities"
	"fmt"
	"strings"

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

func fromProxyToUser(p dbProxyUser) entities.UserEntity {
	return entities.UserEntity{
		ID: p.ID,
		User: entities.User{
			Gmail:               p.Gmail,
			Password:            p.Password,
			Phone:               p.Phone,
			FullName:            p.FullName,
			AllowsAdvertisement: p.AllowsAdvertisement,
			PurchasedRouteIds:   strings.Split(p.PurchasedRouteIds, ","),
		},
	}
}

func fromUserToProxy(u entities.UserEntity) dbProxyUser {
	return dbProxyUser{
		ID:                  u.ID,
		Gmail:               u.Gmail,
		Password:            u.Password,
		Phone:               u.Phone,
		FullName:            u.FullName,
		AllowsAdvertisement: u.AllowsAdvertisement,
		PurchasedRouteIds:   strings.Join(u.PurchasedRouteIds, ","),
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

func (repo sqlUserStorage) Create(user entities.User) (entities.UserEntity, error) {
	entity := entities.UserEntity{User: user, ID: uuid.New().String()}
	proxy := fromUserToProxy(entity)

	result := repo.db.Create(&proxy)
	if result.Error != nil {
		return entities.UserEntity{}, result.Error
	}

	return entity, nil
}

func (repo sqlUserStorage) ReadAll() ([]entities.UserEntity, error) {
	var proxies []dbProxyUser

	result := repo.db.Find(&proxies)
	if result.Error != nil {
		return nil, result.Error
	}

	users := make([]entities.UserEntity, len(proxies))

	for i, proxy := range proxies {
		users[i] = fromProxyToUser(proxy)
	}

	return users, nil
}

func (repo sqlUserStorage) ByID(id string) (entities.UserEntity, error) {
	var proxy dbProxyUser

	result := repo.db.Where("id = ?", id).First(&proxy)
	if result.Error != nil {
		return entities.UserEntity{}, result.Error
	}

	return fromProxyToUser(proxy), nil
}

func (repo sqlUserStorage) Update(userToUpdate entities.UserEntity) error {
	proxy := fromUserToProxy(userToUpdate)

	result := repo.db.Save(&proxy)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo sqlUserStorage) Delete(id string) error {
	result := repo.db.Where("id = ?", id).Delete(&dbProxyUser{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo sqlUserStorage) DropTable() error {
	return repo.db.Migrator().DropTable(&dbProxyUser{})
}
