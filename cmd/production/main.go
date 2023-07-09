package main

import (
	"fmt"
	"time"

	"auth/src/controller"
	"auth/src/notifier"
	"auth/src/security"
	"auth/src/services/passreset"
	"auth/src/services/settings"
	"auth/src/services/signin"
	"auth/src/services/signup"
	appSettings "auth/src/settings"
	"auth/src/storage"
	"auth/src/web"
)

func main() {
	setting := appSettings.NewDotEnvSettings().Load()
	hash := security.Sha256WithSecretFactory(setting.HashSecret)
	futureUserStor := storage.NewRedisTemporaryStorage(setting.RedisAddress, 30*time.Minute, "register")
	resetPassStor := storage.NewRedisTemporaryStorage(setting.RedisAddress, 5*time.Minute, "reset")
	userStorage := storage.NewUserJsonFileStorage("users.json")
	sendFromCreds := notifier.SendFrom{Gmail: setting.Gmail, Password: setting.GmailPassword}

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

	signinService := signin.NewSigninService(userStorage, hash)
	signupService := signup.NewRegistrationService(userStorage, futureUserStor, sendRegisterGmail, security.GenerateCode, hash)
	passwordResetingService := passreset.NewResetPasswordService(userStorage, resetPassStor, sendResetPasswordLetter, hash, security.GenerateCode)
	settingsService := settings.NewSettingsService(userStorage)

	controller := controller.NewController(setting.JwtSecret, signinService, signupService, passwordResetingService, settingsService)

	server := web.SetupServer(controller)

	fmt.Printf("⚡️[server]: Server is running at http://localhost:%d\n", setting.Port)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
