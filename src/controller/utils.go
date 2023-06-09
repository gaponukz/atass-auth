package controller

import (
	"auth/src/entities"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type credentials struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
}

type userCredentialsnDTO struct {
	FullName    string `json:"fullName"`
	Phone       string `json:"phone"`
	RememberHim bool   `json:"rememberHim"`
	credentials
}

type passwordResetConfirmation struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
	Key      string `json:"key"`
}

func getExpirationTime(remember bool) time.Time {
	if remember {
		return time.Now().Add(24 * 11 * time.Hour)
	}

	return time.Now().Add(10 * time.Minute)
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

func getResetPasswordConfirmationFromBody(request *http.Request) (passwordResetConfirmation, error) {
	var creds passwordResetConfirmation

	err := decodeRequestBody(request, &creds)
	if err != nil {
		return passwordResetConfirmation{}, err
	}

	return creds, nil
}

func getGmailFromBody(request *http.Request) (string, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	gmail, ok := data["gmail"].(string)
	if !ok {
		return "", errors.New("gmail field not found or is not a string")
	}

	return gmail, nil
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

func OnlyAuthenticated(router RouterFunc, JwtSecret string) RouterFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		tokenCookie, err := request.Cookie("token")

		if err != nil {
			if err == http.ErrNoCookie {
				responseWriter.WriteHeader(http.StatusUnauthorized)
				return
			}
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		dto, err := getAuthorizedUserDataFromCookie(tokenCookie, JwtSecret)

		if err != nil {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}

		if dto.FullName == "" {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}

		router(responseWriter, request)
	}
}
