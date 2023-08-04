package main

import (
	"fmt"
	"time"

	"auth/src/config"
	"auth/src/consumer"
	"auth/src/controller"
	"auth/src/logger"
	"auth/src/services/passreset"
	"auth/src/services/routes"
	"auth/src/services/session"
	"auth/src/services/settings"
	"auth/src/services/show_routes"
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

	logging := logger.NewConsoleLogger()
	hash := func(s string) string { return s }
	sendRegisterGmail := func(gmail, key string) error { return nil }
	sendResetPasswordLetter := func(gmail, key string) error { return nil }
	generateCode := func() string { return "12345" }

	signinService := logger.NewLogSigninServiceDecorator(signin.NewSigninService(userStorage, hash), logging)
	signupService := logger.NewLogSignupServiceDecorator(signup.NewRegistrationService(userStorage, futureUserStor, sendRegisterGmail, generateCode, hash), logging)
	passwordResetingService := logger.NewLogResetPasswordServiceDecorator(passreset.NewResetPasswordService(userStorage, resetPassStor, sendResetPasswordLetter, hash, generateCode), logging)
	settingsService := logger.NewLogSettingsServiceDecorator(settings.NewSettingsService(userStorage), logging)
	routesService := logger.NewLogAddRouteDecorator(routes.NewRoutesService(userStorage), logging)
	showRoutesService := show_routes.NewShowRoutesService(userStorage)
	sessionService := session.NewSessionService(appSettings.JwtSecret)

	controller := controller.NewController(signinService, signupService, passwordResetingService, settingsService, showRoutesService, sessionService)
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
