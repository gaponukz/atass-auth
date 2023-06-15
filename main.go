package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"auth/src/controller"
	"auth/src/notifier"
	"auth/src/registration"
	"auth/src/resetPassword"
	"auth/src/settings"
	"auth/src/storage"
)

type RouterFunc = func(rw http.ResponseWriter, r *http.Request)

func RequiredMethod(router RouterFunc, required string) RouterFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == required {
			router(responseWriter, request)

		} else {
			http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		log.Printf("%s %s?%s", request.Method, request.URL.Path, request.URL.RawQuery)
		next.ServeHTTP(responseWriter, request)
	})
}

func main() {
	httpRoute := http.NewServeMux()
	settingExporter := settings.DotEnvSettings{}
	settings := settingExporter.Load()
	userStorage := storage.UserJsonFileStorage{
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
			FutureUserStorage: storage.RedisTemporaryStorage(30*time.Minute, "register"),
			Notify: func(gmail, key string) error {
				return sendGmail(gmail, "Confirm Email", notifier.GenerateConfirmCodeLetter(key))
			},
		},
		ResetPasswordService: resetPassword.ResetPasswordService{
			TemporaryStorage: storage.RedisTemporaryStorage(5*time.Minute, "reset"),
			UserStorage:      userStorage,
			Notify: func(gmail, key string) error {
				return sendGmail(gmail, "Confirm password resetting", notifier.GenerateConfirmCodeLetter(key))
			},
		},
	}

	httpRoute.HandleFunc("/signup", RequiredMethod(routerService.Signup, http.MethodPost))
	httpRoute.HandleFunc("/confirmRegistration", RequiredMethod(routerService.ConfirmRegistration, http.MethodPost))
	httpRoute.HandleFunc("/signin", RequiredMethod(routerService.Signin, http.MethodPost))
	httpRoute.HandleFunc("/resetPassword", RequiredMethod(routerService.ResetPassword, http.MethodPost))
	httpRoute.HandleFunc("/confirmResetPassword", RequiredMethod(routerService.ConfirmResetPassword, http.MethodPost))
	httpRoute.HandleFunc("/logout", routerService.Logout)
	httpRoute.HandleFunc("/refresh", routerService.Refresh)

	httpRoute.HandleFunc("/getUserInfo", routerService.GetFullUserInfo)
	httpRoute.HandleFunc("/subscribeUserToTheRoute", RequiredMethod(routerService.SubscribeToTheRoute, http.MethodPost))

	server := http.Server{
		Addr:    ":8080",
		Handler: LoggingMiddleware(httpRoute),
	}

	fmt.Println("⚡️[server]: Server is running at http://localhost:8080")
	server.ListenAndServe()
}
