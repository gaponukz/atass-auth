package main

import (
	"auth/src/controller"
	"auth/src/registration"
	"auth/src/settings"
	"auth/src/storage"
	"fmt"
	"net/http"
)

func main() {
	httpRoute := http.NewServeMux()
	settingExporter := settings.MemorySettingsExporter{}
	settings, _ := settingExporter.Load()
	userStorage := &storage.UserMemoryStorage{}

	routerService := controller.Controller{
		Storage:  userStorage,
		Settings: settings,
		RegistrationService: registration.RegistrationService{
			UserStorage:       userStorage,
			FutureUserStorage: storage.NewFutureUserMemoryStorage(),
			Notify: func(gmail, key string) error {
				fmt.Printf("sent gmail notification for %s with key: %s", gmail, key)
				return nil
			},
		},
	}

	httpRoute.HandleFunc("/signup", routerService.Signup)
	httpRoute.HandleFunc("/confirm", routerService.ConfirmRegistration)
	httpRoute.HandleFunc("/signin", routerService.Signin)
	httpRoute.HandleFunc("/welcome", routerService.Welcome)
	httpRoute.HandleFunc("/refresh", routerService.Refresh)
	httpRoute.HandleFunc("/logout", routerService.Logout)

	server := http.Server{
		Addr:    ":8080",
		Handler: httpRoute,
	}

	fmt.Println("⚡️[server]: Server is running at http://localhost:8080")
	server.ListenAndServe()
}
