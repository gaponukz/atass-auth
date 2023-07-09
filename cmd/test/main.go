package main

import (
	"fmt"
	"time"

	"auth/src/controller"
	"auth/src/security"
	"auth/src/services/passreset"
	"auth/src/services/settings"
	"auth/src/services/signin"
	"auth/src/services/signup"
	"auth/src/storage"
	"auth/src/web"
)

func main() {
	futureUserStor := storage.NewRedisTemporaryStorage("localhost:6379", 1*time.Minute, "register")
	resetPassStor := storage.NewRedisTemporaryStorage("localhost:6379", 1*time.Minute, "reset")
	userStorage := storage.NewUserJsonFileStorage("users.json")

	hash := func(s string) string { return s }
	sendRegisterGmail := func(gmail, key string) error { return nil }
	sendResetPasswordLetter := func(gmail, key string) error { return nil }

	signinService := signin.NewSigninService(userStorage, hash)
	signupService := signup.NewRegistrationService(userStorage, futureUserStor, sendRegisterGmail, security.GenerateCode, hash)
	passwordResetingService := passreset.NewResetPasswordService(userStorage, resetPassStor, sendResetPasswordLetter, hash, security.GenerateCode)
	settingsService := settings.NewSettingsService(userStorage)

	controller := controller.NewController("", signinService, signupService, passwordResetingService, settingsService)
	server := web.SetupTestServer(controller)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
