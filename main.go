package main

import (
	"fmt"
	"net/http"
	"time"

	"auth/src/controller"
	"auth/src/notifier"
	"auth/src/registration"
	"auth/src/resetPassword"
	"auth/src/settings"
	"auth/src/storage"
)

func main() {
	httpRoute := http.NewServeMux()
	settingExporter := settings.DotEnvSettings{}
	settings := settingExporter.Load()
	userStorage := &storage.UserJsonFileStorage{
		FilePath: "users.json",
	}

	sendGmail := notifier.SendEmailNoificationFactory(notifier.SendFrom{
		Gmail:    settings.Gmail,
		Password: settings.GmailPassword,
	})

	routerService := controller.Controller{
		Storage:  userStorage,
		Settings: settings,
		RegistrationService: registration.RegistrationService{
			UserStorage:       userStorage,
			FutureUserStorage: storage.NewFutureUserRedisStorage(30 * time.Minute),
			Notify: func(gmail, key string) error {
				return sendGmail(gmail, "Confirm Email", notifier.GenerateConfirmCodeLetter(key))
			},
		},
		ResetPasswordService: resetPassword.ResetPasswordService{
			TemporaryStorage: storage.NewPasswordResetStorage(),
			UserStorage:      userStorage,
			Notify: func(gmail, key string) error {
				return sendGmail(gmail, "Confirm password resetting", notifier.GenerateConfirmCodeLetter(key))
			},
		},
	}

	httpRoute.HandleFunc("/signup", controller.RequiredMethod(routerService.Signup, http.MethodPost))
	httpRoute.HandleFunc("/confirmRegistration", controller.RequiredMethod(routerService.ConfirmRegistration, http.MethodPost))
	httpRoute.HandleFunc("/signin", controller.RequiredMethod(routerService.Signin, http.MethodPost))
	httpRoute.HandleFunc("/resetPassword", controller.RequiredMethod(routerService.ResetPassword, http.MethodPost))
	httpRoute.HandleFunc("/confirmResetPassword", controller.RequiredMethod(routerService.ConfirmResetPassword, http.MethodPost))
	httpRoute.HandleFunc("/logout", routerService.Logout)
	httpRoute.HandleFunc("/refresh", routerService.Refresh)

	httpRoute.HandleFunc("/getUserInfo", routerService.GetFullUserInfo)
	httpRoute.HandleFunc("/subscribeUserToTheRoute", controller.RequiredMethod(routerService.SubscribeToTheRoute, http.MethodPost))

	server := http.Server{
		Addr:    ":8080",
		Handler: controller.LoggingMiddleware(httpRoute),
	}

	fmt.Println("⚡️[server]: Server is running at http://localhost:8080")
	server.ListenAndServe()
}
