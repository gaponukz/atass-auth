package main

import (
	"fmt"
	"os"
	"time"

	"auth/src/controller"
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
	err := os.WriteFile(databaseFilename, []byte("[]"), 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		os.Remove(databaseFilename)
	}()

	futureUserStor := storage.NewRedisTemporaryStorage("localhost:6379", 1*time.Minute, "register")
	resetPassStor := storage.NewRedisTemporaryStorage("localhost:6379", 1*time.Minute, "reset")
	userStorage := storage.NewUserJsonFileStorage(databaseFilename)

	hash := func(s string) string { return s }
	sendRegisterGmail := func(gmail, key string) error { return nil }
	sendResetPasswordLetter := func(gmail, key string) error { return nil }
	generateCode := func() string { return "12345" }

	signinService := signin.NewSigninService(userStorage, hash)
	signupService := signup.NewRegistrationService(userStorage, futureUserStor, sendRegisterGmail, generateCode, hash)
	passwordResetingService := passreset.NewResetPasswordService(userStorage, resetPassStor, sendResetPasswordLetter, hash, generateCode)
	settingsService := settings.NewSettingsService(userStorage)
	showRoutesService := show_routes.NewShowRoutesService(userStorage)

	controller := controller.NewController("", signinService, signupService, passwordResetingService, settingsService, showRoutesService)
	server := web.SetupTestServer(controller)

	go func() {
		time.Sleep(time.Second * 3)
		users, _ := userStorage.ReadAll()

		if len(users) != 0 {
			user := users[0]

			user.PurchasedRouteIds = append(user.PurchasedRouteIds, "12-34-5")
			_ = userStorage.Update(user)
		}
	}()

	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
