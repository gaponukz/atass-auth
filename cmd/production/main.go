package main

import (
	"fmt"
	"time"

	"auth/src/controller"
	"auth/src/notifier"
	"auth/src/registration"
	"auth/src/resetPassword"
	"auth/src/security"
	"auth/src/settings"
	"auth/src/storage"
	"auth/src/web"
)

func main() {
	settings := settings.NewDotEnvSettings().Load()
	futureUserStor := storage.NewRedisTemporaryStorage(30*time.Minute, "register")
	resetPassStor := storage.NewRedisTemporaryStorage(5*time.Minute, "reset")
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
			Storage:              userStorage,
			Settings:             settings,
			RegistrationService:  registration.NewRegistrationService(userStorage, futureUserStor, sendRegisterGmail, security.GenerateCode),
			ResetPasswordService: resetPassword.NewResetPasswordService(userStorage, resetPassStor, sendResetPasswordLetter, security.GenerateCode),
		},
	)

	fmt.Printf("⚡️[server]: Server is running at http://localhost:%d", settings.Port)
	server.ListenAndServe()
}
