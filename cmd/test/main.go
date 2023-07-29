package main

import (
	"fmt"
	"time"

	"auth/src/config"
	"auth/src/consumer"
	"auth/src/controller"
	"auth/src/services/passreset"
	"auth/src/services/routes"
	"auth/src/services/settings"
	"auth/src/services/signin"
	"auth/src/services/signup"
	"auth/src/storage"
	"auth/src/web"
)

func main() {
	appSettings := config.NewDotEnvSettings().Load()
	futureUserStor := storage.NewRedisTemporaryStorage(appSettings.RedisAddress, 1*time.Minute, "register")
	resetPassStor := storage.NewRedisTemporaryStorage(appSettings.RedisAddress, 1*time.Minute, "reset")
	userStorage := storage.NewUserJsonFileStorage("users.json")

	hash := func(s string) string { return s }
	sendRegisterGmail := func(gmail, key string) error { return nil }
	sendResetPasswordLetter := func(gmail, key string) error { return nil }
	generateCode := func() string { return "12345" }

	signinService := signin.NewSigninService(userStorage, hash)
	signupService := signup.NewRegistrationService(userStorage, futureUserStor, sendRegisterGmail, generateCode, hash)
	passwordResetingService := passreset.NewResetPasswordService(userStorage, resetPassStor, sendResetPasswordLetter, hash, generateCode)
	settingsService := settings.NewSettingsService(userStorage)
	routesService := routes.NewRoutesService(userStorage)

	controller := controller.NewController("", signinService, signupService, passwordResetingService, settingsService)
	server := web.SetupTestServer(controller)

	routesEventsListener, err := consumer.NewRoutesEventsListener(routesService, appSettings.RabbitUrl)
	if err != nil {
		panic(err.Error())
	}

	defer routesEventsListener.Close()

	go routesEventsListener.Listen()

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
