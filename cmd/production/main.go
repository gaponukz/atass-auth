package main

import (
	"fmt"
	"time"

	"auth/src/controller"
	"auth/src/notifier"
	"auth/src/password_reseting"
	"auth/src/registration"
	"auth/src/security"
	"auth/src/settings"
	"auth/src/storage"
	"auth/src/web"
)

func main() {
	settings := settings.NewDotEnvSettings().Load()
	hash := security.Sha256WithSecretFactory(settings.HashSecret)
	futureUserStor := storage.NewRedisTemporaryStorage(settings.RedisAddress, 30*time.Minute, "register")
	resetPassStor := storage.NewRedisTemporaryStorage(settings.RedisAddress, 5*time.Minute, "reset")
	userStorage := storage.NewUserJsonFileStorage("users.json")
	sendFromCreds := notifier.SendFrom{Gmail: settings.Gmail, Password: settings.GmailPassword}

	sendRegisterGmail := notifier.SendEmailNoificationFactory(
		sendFromCreds,
		"Confirm your registration",
		"letters/confirmRegistration.html",
	)

	sendResetPasswordLetter := notifier.SendEmailNoificationFactory(
		sendFromCreds,
		"Confirm your password reseting",
		"letters/resetPasswors.html",
	)

	server := web.SetupServer(
		controller.Controller{
			Storage:  userStorage,
			Settings: settings,
			RegistrationService: registration.NewRegistrationService(
				userStorage,
				futureUserStor,
				sendRegisterGmail,
				security.GenerateCode,
				hash,
			),
			ResetPasswordService: password_reseting.NewResetPasswordService(
				userStorage,
				resetPassStor,
				sendResetPasswordLetter,
				hash,
				security.GenerateCode,
			),
		},
	)

	fmt.Printf("⚡️[server]: Server is running at http://localhost:%d\n", settings.Port)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
