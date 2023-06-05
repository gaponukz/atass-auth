package controller

import (
	"auth/src/entities"
	"encoding/json"
	"log"
	"net/http"
)

type credentials struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
}

type userCredentialsnDTO struct {
	FullName string `json:"fullName"`
	credentials
}

func decodeRequestBody(request *http.Request, data interface{}) error {
	err := json.NewDecoder(request.Body).Decode(data)

	if err != nil {
		return err
	}

	return nil
}

func getUserCredentialsFromBody(request *http.Request) (userCredentialsnDTO, error) {
	var creds userCredentialsnDTO

	err := decodeRequestBody(request, &creds)
	if err != nil {
		return userCredentialsnDTO{}, err
	}

	return creds, nil
}

func getGmailConfirmationFromBody(request *http.Request) (entities.FutureUser, error) {
	var creds entities.FutureUser

	err := decodeRequestBody(request, &creds)
	if err != nil {
		return entities.FutureUser{}, err
	}

	return creds, nil
}

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
