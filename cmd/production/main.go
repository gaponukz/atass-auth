package main

import (
	"auth/src/application/usecases/passreset"
	"auth/src/application/usecases/session"
	"auth/src/application/usecases/settings"
	"auth/src/application/usecases/show_routes"
	"auth/src/application/usecases/signin"
	"auth/src/application/usecases/signup"
	"auth/src/infr/config"
	"auth/src/infr/gmail_notifier"
	"auth/src/infr/logger"
	"auth/src/infr/security"
	"auth/src/infr/storage"
	"auth/src/interface/controller"
	"fmt"
	"time"
)

func main() {
	setting := config.NewDotEnvSettings().Load()

	userStorage, err := storage.NewMySQLUserStorage(storage.MySQLCredentials{
		Host:     setting.MysqlHost,
		User:     setting.MysqlUser,
		Password: setting.MysqlPassword,
		Dbname:   setting.MysqlDbname,
		Port:     setting.MysqlPort,
	})
	if err != nil {
		panic(err.Error())
	}

	hash := security.Sha256WithSecretFactory(setting.HashSecret)
	futureUserStor := storage.NewRedisTemporaryStorage(setting.RedisAddress, 30*time.Minute, "register")
	resetPassStor := storage.NewRedisTemporaryStorage(setting.RedisAddress, 5*time.Minute, "reset")
	sendFromCreds := gmail_notifier.GmailCreds{Gmail: setting.Gmail, Password: setting.GmailPassword}

	signupNotifier := gmail_notifier.NewGmailNotifier(sendFromCreds, gmail_notifier.Letter{
		Title:    "Confirm your registration",
		HtmlPath: "letters/confirmRegistration.html",
	})

	passresetNotifier := gmail_notifier.NewGmailNotifier(sendFromCreds, gmail_notifier.Letter{
		Title:    "Confirm your password reseting",
		HtmlPath: "letters/resetPasswors.html",
	})

	logging := logger.NewConsoleLogger()
	signinService := logger.NewLogSigninServiceDecorator(signin.NewSigninService(userStorage, hash), logging)
	signupService := logger.NewLogSignupServiceDecorator(signup.NewRegistrationService(userStorage, futureUserStor, signupNotifier, security.GenerateCode, hash), logging)
	passwordResetingService := logger.NewLogResetPasswordServiceDecorator(passreset.NewResetPasswordService(userStorage, resetPassStor, passresetNotifier, hash, security.GenerateCode), logging)
	settingsService := logger.NewLogSettingsServiceDecorator(settings.NewSettingsService(userStorage), logging)
	// routesService := logger.NewLogAddRouteDecorator(routes.NewRoutesService(userStorage), logging)
	showRoutesService := show_routes.NewShowRoutesService(userStorage)
	sessionService := session.NewSessionService(setting.JwtSecret)

	contr := controller.NewController(signinService, signupService, passwordResetingService, settingsService, showRoutesService, sessionService)
	// routesEventsListener, err := event_handler.NewRoutesEventsListener(routesService, setting.RabbitUrl)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// defer routesEventsListener.Close()

	// go routesEventsListener.Listen()

	server := controller.SetupTestServer(contr)

	fmt.Printf("⚡️[server]: Server is running at http://localhost:%d\n", setting.Port)

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
