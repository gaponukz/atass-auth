package main

import (
	"auth/src/controller"
	"auth/src/settings"
	"auth/src/storage"
	"fmt"
	"net/http"
)

func main() {
	httpRoute := http.NewServeMux()
	settingExporter := settings.MemorySettingsExporter{}
	settings, _ := settingExporter.Load()

	routerService := controller.Controller{
		Storage:  &storage.UserMemoryStorage{},
		Settings: settings,
	}

	httpRoute.HandleFunc("/signup", routerService.Signup)
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
