package main

import (
	"fmt"
	"time"

	"auth/src/controller"
	"auth/src/services/passreset"
	"auth/src/services/settings"
	"auth/src/services/signin"
	"auth/src/services/signup"
	"auth/src/storage"
	"auth/src/web"
)

func main() {
	creds := storage.PostgresCredentials{
		Host:     "localhost",
		User:     "myuser",
		Password: "mypassword",
		Dbname:   "users",
		Port:     "5432",
		Sslmode:  "disable",
	}
	userStorage, err := storage.NewPostgresUserStorage(creds)
	if err != nil {
		panic(err)
	}

	defer userStorage.DropTable()

	futureUserStor := storage.NewRedisTemporaryStorage("localhost:6379", 1*time.Minute, "register")
	resetPassStor := storage.NewRedisTemporaryStorage("localhost:6379", 1*time.Minute, "reset")

	hash := func(s string) string { return s }
	sendRegisterGmail := func(gmail, key string) error { return nil }
	sendResetPasswordLetter := func(gmail, key string) error { return nil }
	generateCode := func() string { return "12345" }

	signinService := signin.NewSigninService(userStorage, hash)
	signupService := signup.NewRegistrationService(userStorage, futureUserStor, sendRegisterGmail, generateCode, hash)
	passwordResetingService := passreset.NewResetPasswordService(userStorage, resetPassStor, sendResetPasswordLetter, hash, generateCode)
	settingsService := settings.NewSettingsService(userStorage)

	controller := controller.NewController("", signinService, signupService, passwordResetingService, settingsService)
	server := web.SetupTestServer(controller)

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
