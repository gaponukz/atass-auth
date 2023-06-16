package main

import (
	"fmt"
	"time"

	"auth/src/controller"
	"auth/src/registration"
	"auth/src/resetPassword"
	"auth/src/settings"
	"auth/src/storage"
	"auth/src/web"
)

func main() {
	settings := settings.DotEnvSettings{}.Load()
	futureUserStor := storage.RedisTemporaryStorage(30*time.Minute, "register")
	resetPassStor := storage.RedisTemporaryStorage(5*time.Minute, "reset")
	userStorage := storage.UserJsonFileStorage{FilePath: "users.json"}

	server := web.SetupServer(
		controller.Controller{
			Storage:  userStorage,
			Settings: settings,
			RegistrationService: registration.RegistrationService{
				UserStorage:       userStorage,
				FutureUserStorage: futureUserStor,
				Notify:            func(gmail string, key string) error { return nil },
				GenerateCode:      func() string { return "12345" },
			},
			ResetPasswordService: resetPassword.ResetPasswordService{
				TemporaryStorage: resetPassStor,
				UserStorage:      userStorage,
				Notify:           func(gmail string, key string) error { return nil },
				GenerateCode:     func() string { return "12345" },
			},
		},
	)

	fmt.Printf("⚡️[server]: Server is running at http://localhost:%d", settings.Port)
	server.ListenAndServe()
}
