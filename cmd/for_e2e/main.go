package main

import (
	"fmt"
	"os"
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
	err := os.WriteFile(databaseFilename, []byte("[]"), 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		os.Remove(databaseFilename)
	}()

	settings := settings.NewDotEnvSettings().Load()
	futureUserStor := storage.NewRedisTemporaryStorage("localhost:6379", 1*time.Minute, "register")
	resetPassStor := storage.NewRedisTemporaryStorage("localhost:6379", 1*time.Minute, "reset")
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

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
