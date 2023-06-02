package main

import (
	"auth/src/controller"
	"auth/src/storage"
	"fmt"
	"net/http"
)

func main() {
	httpRoute := http.NewServeMux()
	routerService := controller.Controller{
		Storage: &storage.MemoryStorage{},
	}

	httpRoute.HandleFunc("/signin", routerService.Singin)
	httpRoute.HandleFunc("/refresh", routerService.Refresh)
	httpRoute.HandleFunc("/logout", routerService.Logout)
	httpRoute.HandleFunc("/welcome", routerService.Welcome)

	server := http.Server{
		Addr:    ":8080",
		Handler: httpRoute,
	}

	fmt.Println("⚡️[server]: Server is running at http://localhost:8080")
	server.ListenAndServe()
}
