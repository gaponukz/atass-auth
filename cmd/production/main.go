package main

import (
	"auth/src/application/usecases/passreset"
	"auth/src/application/usecases/session"
	"auth/src/application/usecases/settings"
	"auth/src/application/usecases/show_routes"
	"auth/src/application/usecases/signin"
	"auth/src/application/usecases/signup"
	"auth/src/infr/config"
	"auth/src/infr/notifier"
	"auth/src/infr/security"
	"auth/src/infr/storage"
	"auth/src/interface/controller"
	"fmt"
	"time"
)

func main() {
	setting := config.NewDotEnvSettings().Load()
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
	showRoutesService := show_routes.NewShowRoutesService(userStorage)
	sessionService := session.NewSessionService(setting.JwtSecret)

	contr := controller.NewController(signinService, signupService, passwordResetingService, settingsService, showRoutesService, sessionService)

	server := controller.SetupServer(contr)

	fmt.Printf("⚡️[server]: Server is running at http://localhost:%d\n", setting.Port)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
