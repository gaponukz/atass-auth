package main

import (
	"auth/src/entities"
	"auth/src/storage"
	"fmt"
)

func main() {
	creds := storage.PostgresCredentials{
		Host:     "localhost",
		User:     "myuser",
		Password: "mypassword",
		Dbname:   "users",
		Port:     "5432",
		Sslmode:  "disable",
	}
	repo, err := storage.NewPostgresUserStorage(creds)
	if err != nil {
		panic(err)
	}

	defer repo.DropTable()

	newUser := entities.User{
		Gmail:               "example@gmail.com",
		Password:            "secret",
		Phone:               "1234567890",
		FullName:            "John Doe",
		AllowsAdvertisement: true,
		PurchasedRouteIds:   []string{"route1", "route2"},
	}

	_, err = repo.Create(newUser)
	if err != nil {
		panic(err)
	}

	allUsers, err := repo.ReadAll()
	if err != nil {
		panic(err)
	}
	fmt.Println(allUsers)
}
