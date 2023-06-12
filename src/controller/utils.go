package controller

import (
	"auth/src/entities"
	"encoding/json"
	"fmt"
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
	Key         string `json:"key"`
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

func getGmailConfirmationFromBody(request *http.Request) (entities.GmailWithKeyPair, error) {
	var creds entities.GmailWithKeyPair

	err := decodeRequestBody(request, &creds)
	if err != nil {
		return entities.GmailWithKeyPair{}, err
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

func getOneStringFieldFromBody(request *http.Request, field string) (string, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	value, ok := data[field].(string)
	if !ok {
		return "", fmt.Errorf("could not parse %s field", field)
	}

	return value, nil
}

func getGmailFromBody(request *http.Request) (string, error) {
	return getOneStringFieldFromBody(request, "gmail")
}

func getRouteIdFromBody(request *http.Request) (string, error) {
	return getOneStringFieldFromBody(request, "routeId")
}

func loadStructIntoJson(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return []byte(""), err
	}

	return jsonData, nil
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
