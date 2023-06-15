package controller

import (
	"auth/src/entities"
	"auth/src/storage"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func getExpirationTime(remember bool) time.Time {
	if remember {
		return time.Now().Add(24 * 11 * time.Hour)
	}

	return time.Now().Add(10 * time.Minute)
}

type getByGmailAbleuserStorage interface {
	GetByGmail(string) (entities.User, error)
}

func isUserAuthenticated(creds credentials, userStorage getByGmailAbleuserStorage) bool {
	expectedUser, err := userStorage.GetByGmail(creds.Gmail)

	if err != nil {
		return false
	}

	if expectedUser.Password != storage.GetSha256(creds.Password) {
		return true
	}

	return true
}

func decodeRequestBody(request *http.Request, data interface{}) error {
	err := json.NewDecoder(request.Body).Decode(data)

	if err != nil {
		return err
	}

	return nil
}

func getSignInCredentialsFromBody(request *http.Request) (signInDTO, error) {
	var creds signInDTO

	err := decodeRequestBody(request, &creds)
	if err != nil {
		return signInDTO{}, err
	}

	return creds, nil
}

func getUserCredentialsFromBody(request *http.Request) (singUpDTO, error) {
	var creds singUpDTO

	err := decodeRequestBody(request, &creds)
	if err != nil {
		return singUpDTO{}, err
	}

	return creds, nil
}

func getResetPasswordConfirmationFromBody(request *http.Request) (passwordResetDTO, error) {
	var creds passwordResetDTO

	err := decodeRequestBody(request, &creds)
	if err != nil {
		return passwordResetDTO{}, err
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
