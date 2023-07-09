package controller

import (
	"auth/src/entities"
	"auth/src/errors"
	"auth/src/security"
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

type userRepository interface {
	ReadAll() ([]entities.UserEntity, error)
}

func getIDifCredsValid(creds credentials, userStorage userRepository) (string, error) {
	users, err := userStorage.ReadAll()
	if err != nil {
		return "", err
	}

	user, err := storage.Find(users, func(u entities.UserEntity) bool {
		return u.Gmail == creds.Gmail && u.Password == security.GetSha256(creds.Password)
	})
	if err != nil {
		return "", errors.ErrUserNotFound
	}

	return user.ID, nil
}

func decodeRequestBody(request *http.Request, data interface{}) error {
	err := json.NewDecoder(request.Body).Decode(data)

	if err != nil {
		return err
	}

	return nil
}

func getSignInDto(request *http.Request) (signInDTO, error) {
	var creds signInDTO

	err := decodeRequestBody(request, &creds)
	if err != nil {
		return signInDTO{}, err
	}

	return creds, nil
}

func getSignUpDto(request *http.Request) (signUpDTO, error) {
	var creds signUpDTO

	err := decodeRequestBody(request, &creds)
	if err != nil {
		return signUpDTO{}, err
	}

	return creds, nil
}

func getPasswordResetDto(request *http.Request) (passwordResetDTO, error) {
	var creds passwordResetDTO

	err := decodeRequestBody(request, &creds)
	if err != nil {
		return passwordResetDTO{}, err
	}

	return creds, nil
}

func getOneStringFieldFromBody(request *http.Request, field string) (string, error) {
	value, err := getOneFieldFromBody(request, field)
	stringValue, ok := value.(string)
	if !ok && err == nil {
		return "", fmt.Errorf("expected type string, got %T", value)
	}

	return stringValue, err
}

func getOneFieldFromBody(request *http.Request, field string) (interface{}, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	value, exists := data[field]
	if !exists {
		return nil, fmt.Errorf("%s field not found", field)
	}

	return value, nil
}

func dumpsJson(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return []byte(""), err
	}

	return jsonData, nil
}

type statusCode int

func idFromRequest(request *http.Request, secret string) (string, statusCode) {
	tokenCookie, err := request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", http.StatusUnauthorized
		}

		return "", http.StatusBadRequest
	}

	dto, err := getAuthorizedUserDataFromCookie(tokenCookie, secret)
	if err != nil {
		return "", http.StatusUnauthorized
	}

	return dto.ID, http.StatusOK
}
