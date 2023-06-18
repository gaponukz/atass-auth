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
	settings := settings.NewDotEnvSettings().Load()
	futureUserStor := storage.NewRedisTemporaryStorage(30*time.Minute, "register")
	resetPassStor := storage.NewRedisTemporaryStorage(5*time.Minute, "reset")
	userStorage := storage.NewUserJsonFileStorage("users.json")

	server := web.SetupTestServer(
		controller.Controller{
			Storage:  userStorage,
			Settings: settings,
			RegistrationService: registration.NewRegistrationService(
				userStorage,
				futureUserStor,
				func(gmail string, key string) error { return nil },
				func() string { return "12345" },
			),
			ResetPasswordService: resetPassword.NewResetPasswordService(
				userStorage,
				resetPassStor,
				func(gmail string, key string) error { return nil },
				func() string { return "12345" },
			),
		},
	)

	fmt.Printf("⚡️[server]: Server is running at http://localhost:%d", settings.Port)
	server.ListenAndServe()
}
