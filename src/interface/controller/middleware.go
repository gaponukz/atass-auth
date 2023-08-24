package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type routerFunc = func(rw http.ResponseWriter, r *http.Request)

func requiredMethod(router routerFunc, required string) routerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method == required {
			router(responseWriter, request)

		} else {
			http.Error(responseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		log.Printf("%s %s?%s", request.Method, request.URL.Path, request.URL.RawQuery)
		next.ServeHTTP(responseWriter, request)
	})
}

func enableCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Accept, Content-Type, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "1728000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func getMuxFromController(c *Controller) *http.ServeMux {
	httpRoute := http.NewServeMux()

	httpRoute.HandleFunc("/api/auth/signup", requiredMethod(c.Signup, http.MethodPost))
	httpRoute.HandleFunc("/api/auth/confirmRegistration", requiredMethod(c.ConfirmRegistration, http.MethodPost))
	httpRoute.HandleFunc("/api/auth/signin", requiredMethod(c.Signin, http.MethodPost))
	httpRoute.HandleFunc("/api/auth/resetPassword", requiredMethod(c.ResetPassword, http.MethodPost))
	httpRoute.HandleFunc("/api/auth/confirmResetPassword", requiredMethod(c.ConfirmResetPassword, http.MethodPost))
	httpRoute.HandleFunc("/api/auth/cancelPasswordResetting", requiredMethod(c.CancelPasswordResetting, http.MethodPost))
	httpRoute.HandleFunc("/api/auth/logout", c.Logout)
	httpRoute.HandleFunc("/api/auth/refresh", c.Refresh)

	httpRoute.HandleFunc("/api/auth/getUserRoutes", c.ShowUserRoutes)
	httpRoute.HandleFunc("/api/auth/getUserInfo", c.GetUserInfo)
	httpRoute.HandleFunc("/api/auth/updateUserInfo", requiredMethod(c.UpdateUserInfo, http.MethodPost))

	return httpRoute
}

func SetupServer(c *Controller) *http.Server {
	handler := getMuxFromController(c)

	return &http.Server{
		Addr:              ":8080",
		Handler:           loggingMiddleware(handler),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func SetupTestServer(c *Controller) *http.Server {
	handler := getMuxFromController(c)

	fmt.Println("Warning: this is test server, please do not use it in production.")
	return &http.Server{
		Addr:              ":8080",
		Handler:           enableCORS(loggingMiddleware(handler)),
		ReadHeaderTimeout: 2 * time.Second,
	}
}
