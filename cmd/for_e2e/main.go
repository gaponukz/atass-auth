package main

import (
	"auth/src/application/usecases/passreset"
	"auth/src/application/usecases/session"
	"auth/src/application/usecases/settings"
	"auth/src/application/usecases/show_routes"
	"auth/src/application/usecases/signin"
	"auth/src/application/usecases/signup"
	"auth/src/domain/entities"
	"auth/src/infr/logger"
	"auth/src/infr/storage"
	"auth/src/interface/controller"
	"fmt"
	"os"
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
	generateCode := func() string { return "12345" }

	signinService := logger.NewLogSigninServiceDecorator(signin.NewSigninService(userStorage, hash), logging)
	signupService := logger.NewLogSignupServiceDecorator(signup.NewRegistrationService(userStorage, futureUserStor, notifierMock{}, generateCode, hash), logging)
	passwordResetingService := logger.NewLogResetPasswordServiceDecorator(passreset.NewResetPasswordService(userStorage, resetPassStor, notifierMock{}, hash, generateCode), logging)
	settingsService := logger.NewLogSettingsServiceDecorator(settings.NewSettingsService(userStorage), logging)
	showRoutesService := show_routes.NewShowRoutesService(userStorage)
	sessionService := session.NewSessionService("guaghf79gf")

	contr := controller.NewController(signinService, signupService, passwordResetingService, settingsService, showRoutesService, sessionService)
	server := controller.SetupServer(contr)

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
