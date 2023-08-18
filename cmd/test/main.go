package main

import (
	"auth/src/application/usecases/passreset"
	"auth/src/application/usecases/routes"
	"auth/src/application/usecases/session"
	"auth/src/application/usecases/settings"
	"auth/src/application/usecases/show_routes"
	"auth/src/application/usecases/signin"
	"auth/src/application/usecases/signup"
	"auth/src/domain/entities"
	"auth/src/infr/config"
	"auth/src/infr/logger"
	"auth/src/infr/storage"
	"auth/src/interface/controller"
	"auth/src/interface/event_handler"
	"fmt"
	"time"
)

type notifierMock struct{}

func (m notifierMock) Notify(to, code string) error {
	return nil
}

func (m notifierMock) NotifyUser(to entities.User, code string) error {
	return nil
}

func main() {
	appSettings := config.NewDotEnvSettings().Load()
	futureUserStor := storage.NewRedisTemporaryStorage(appSettings.RedisAddress, 1*time.Minute, "register")
	resetPassStor := storage.NewRedisTemporaryStorage(appSettings.RedisAddress, 1*time.Minute, "reset")
	userStorage := storage.NewUserJsonFileStorage("users.json")

	logging := logger.NewConsoleLogger()
	hash := func(s string) string { return s }
	generateCode := func() string { return "12345" }

	signinService := logger.NewLogSigninServiceDecorator(signin.NewSigninService(userStorage, hash), logging)
	signupService := logger.NewLogSignupServiceDecorator(signup.NewRegistrationService(userStorage, futureUserStor, notifierMock{}, generateCode, hash), logging)
	passwordResetingService := logger.NewLogResetPasswordServiceDecorator(passreset.NewResetPasswordService(userStorage, resetPassStor, notifierMock{}, hash, generateCode), logging)
	settingsService := logger.NewLogSettingsServiceDecorator(settings.NewSettingsService(userStorage), logging)
	addRouteService := logger.NewLogAddRouteDecorator(routes.NewAddRouteService(userStorage), logging)
	deleteRouteService := logger.NewLogDeleteRouteDecorator(routes.NewDeleteRouteService(userStorage), logging)
	showRoutesService := show_routes.NewShowRoutesService(userStorage)
	sessionService := session.NewSessionService(appSettings.JwtSecret)

	contr := controller.NewController(signinService, signupService, passwordResetingService, settingsService, showRoutesService, sessionService)
	server := controller.SetupTestServer(contr)

	routesEventsListener, err := event_handler.NewRoutesEventsListener(addRouteService, deleteRouteService, appSettings.RabbitUrl)
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
