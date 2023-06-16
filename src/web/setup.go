package web

import (
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

type controller interface {
	Signin(rw http.ResponseWriter, r *http.Request)
	Signup(rw http.ResponseWriter, r *http.Request)
	ConfirmRegistration(rw http.ResponseWriter, r *http.Request)
	Refresh(rw http.ResponseWriter, r *http.Request)
	Logout(rw http.ResponseWriter, r *http.Request)
	ResetPassword(rw http.ResponseWriter, r *http.Request)
	ConfirmResetPassword(rw http.ResponseWriter, r *http.Request)
	GetFullUserInfo(rw http.ResponseWriter, r *http.Request)
	SubscribeToTheRoute(rw http.ResponseWriter, r *http.Request)
}

func SetupServer(c controller) *http.Server {
	httpRoute := http.NewServeMux()

	httpRoute.HandleFunc("/signup", requiredMethod(c.Signup, http.MethodPost))
	httpRoute.HandleFunc("/confirmRegistration", requiredMethod(c.ConfirmRegistration, http.MethodPost))
	httpRoute.HandleFunc("/signin", requiredMethod(c.Signin, http.MethodPost))
	httpRoute.HandleFunc("/resetPassword", requiredMethod(c.ResetPassword, http.MethodPost))
	httpRoute.HandleFunc("/confirmResetPassword", requiredMethod(c.ConfirmResetPassword, http.MethodPost))
	httpRoute.HandleFunc("/logout", c.Logout)
	httpRoute.HandleFunc("/refresh", c.Refresh)

	httpRoute.HandleFunc("/getUserInfo", c.GetFullUserInfo)
	httpRoute.HandleFunc("/subscribeUserToTheRoute", requiredMethod(c.SubscribeToTheRoute, http.MethodPost))

	return &http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(httpRoute),
	}
}
