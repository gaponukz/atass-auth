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
	settings := settings.DotEnvSettings{}.Load()
	futureUserStor := storage.RedisTemporaryStorage(30*time.Minute, "register")
	resetPassStor := storage.RedisTemporaryStorage(5*time.Minute, "reset")
	userStorage := storage.UserJsonFileStorage{FilePath: "users.json"}
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
			RegistrationService: registration.RegistrationService{
				UserStorage:       userStorage,
				FutureUserStorage: futureUserStor,
				Notify:            sendRegisterGmail,
				GenerateCode:      security.GenerateCode,
			},
			ResetPasswordService: resetPassword.ResetPasswordService{
				TemporaryStorage: resetPassStor,
				UserStorage:      userStorage,
				Notify:           sendResetPasswordLetter,
				GenerateCode:     security.GenerateCode,
			},
		},
	)

	fmt.Println("⚡️[server]: Server is running at http://localhost:8080")
	server.ListenAndServe()
}
