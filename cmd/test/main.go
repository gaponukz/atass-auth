package main

import (
	"fmt"
	"time"

	"auth/src/controller"
	"auth/src/password_reseting"
	"auth/src/registration"
	"auth/src/settings"
	"auth/src/storage"
	"auth/src/web"
)

func main() {
	databaseFilename := "test.json"
	settings := settings.NewDotEnvSettings().Load()
	futureUserStor := storage.NewRedisTemporaryStorage(settings.RedisAddress, 30*time.Minute, "register")
	resetPassStor := storage.NewRedisTemporaryStorage(settings.RedisAddress, 5*time.Minute, "reset")
	userStorage := storage.NewUserJsonFileStorage(databaseFilename)

	server := web.SetupTestServer(
		controller.Controller{
			Storage:  userStorage,
			Settings: settings,
			RegistrationService: registration.NewRegistrationService(
				userStorage,
				futureUserStor,
				func(gmail string, key string) error { return nil },
				func() string { return "12345" },
				func(s string) string { return s },
			),
			ResetPasswordService: password_reseting.NewResetPasswordService(
				userStorage,
				resetPassStor,
				func(gmail string, key string) error { return nil },
				func(s string) string { return s },
				func() string { return "12345" },
			),
		},
	)

	fmt.Printf("⚡️[redis]: is running at: %s\n", settings.RedisAddress)
	fmt.Println("⚡️[server]: is running at http://localhost:8080")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
