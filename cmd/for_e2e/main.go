package main

import (
	"fmt"
	"os"
	"time"

	"auth/src/controller"
	"auth/src/logger"
	"auth/src/services/passreset"
	"auth/src/services/settings"
	"auth/src/services/show_routes"
	"auth/src/services/signin"
	"auth/src/services/signup"
	"auth/src/storage"
	"auth/src/web"
)

func main() {
	databaseFilename := "test.json"
	err := os.WriteFile(databaseFilename, []byte("[]"), 0600)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.Remove(databaseFilename)
	}()

	futureUserStor := storage.NewRedisTemporaryStorage("localhost:6379", 1*time.Minute, "register")
	resetPassStor := storage.NewRedisTemporaryStorage("localhost:6379", 1*time.Minute, "reset")
	userStorage := storage.NewUserJsonFileStorage(databaseFilename)

	logging := logger.NewConsoleLogger()
	hash := func(s string) string { return s }
	sendRegisterGmail := func(gmail, key string) error { return nil }
	sendResetPasswordLetter := func(gmail, key string) error { return nil }
	generateCode := func() string { return "12345" }

	signinService := logger.NewLogSigninServiceDecorator(signin.NewSigninService(userStorage, hash), logging)
	signupService := logger.NewLogSignupServiceDecorator(signup.NewRegistrationService(userStorage, futureUserStor, sendRegisterGmail, generateCode, hash), logging)
	passwordResetingService := logger.NewLogResetPasswordServiceDecorator(passreset.NewResetPasswordService(userStorage, resetPassStor, sendResetPasswordLetter, hash, generateCode), logging)
	settingsService := logger.NewLogSettingsServiceDecorator(settings.NewSettingsService(userStorage), logging)
	showRoutesService := show_routes.NewShowRoutesService(userStorage)

	controller := controller.NewController("", signinService, signupService, passwordResetingService, settingsService, showRoutesService)
	server := web.SetupTestServer(controller)

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
