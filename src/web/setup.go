package web

import (
	"auth/src/controller"
	"fmt"
	"log"
	"net/http"
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

func getMuxFromController(c *controller.Controller) *http.ServeMux {
	httpRoute := http.NewServeMux()

	httpRoute.HandleFunc("/signup", requiredMethod(c.Signup, http.MethodPost))
	httpRoute.HandleFunc("/confirmRegistration", requiredMethod(c.ConfirmRegistration, http.MethodPost))
	httpRoute.HandleFunc("/signin", requiredMethod(c.Signin, http.MethodPost))
	httpRoute.HandleFunc("/resetPassword", requiredMethod(c.ResetPassword, http.MethodPost))
	httpRoute.HandleFunc("/confirmResetPassword", requiredMethod(c.ConfirmResetPassword, http.MethodPost))
	httpRoute.HandleFunc("/cancelPasswordResetting", requiredMethod(c.CancelPasswordResetting, http.MethodPost))
	httpRoute.HandleFunc("/logout", c.Logout)
	httpRoute.HandleFunc("/refresh", c.Refresh)

	httpRoute.HandleFunc("/getUserRoutes", c.ShowUserRoutes)
	httpRoute.HandleFunc("/getUserInfo", c.GetUserInfo)
	httpRoute.HandleFunc("/updateUserInfo", requiredMethod(c.UpdateUserInfo, http.MethodPost))

	return httpRoute
}

func SetupServer(c *controller.Controller) *http.Server {
	handler := getMuxFromController(c)

	return &http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(handler),
	}
}

func SetupTestServer(c *controller.Controller) *http.Server {
	handler := getMuxFromController(c)

	fmt.Println("Warning: this is test server, please do not use it in production.")
	return &http.Server{
		Addr:    ":8080",
		Handler: enableCORS(loggingMiddleware(handler)),
	}
}
