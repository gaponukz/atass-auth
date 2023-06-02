package controller

import (
	"auth/src/entities"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type credentials struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
}

type claims struct {
	FullName string `json:"fullName"`
	jwt.RegisteredClaims
}

type GetByGmailAbleStorage interface {
	GetByGmail(string) (entities.User, error)
}

func getRegisteredUserFromRequestBody(request *http.Request, storage GetByGmailAbleStorage) (entities.User, error) {
	var creds credentials

	err := json.NewDecoder(request.Body).Decode(&creds)
	if err != nil {
		return entities.User{}, err
	}

	expectedUser, err := storage.GetByGmail(creds.Gmail)

	if err != nil {
		return entities.User{}, err
	}

	if expectedUser.Password != creds.Password {
		return entities.User{}, fmt.Errorf("unauthorized")
	}

	return expectedUser, nil
}

func getTemporaryToken(userFullName, jwtSecret string) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &claims{
		FullName: userFullName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}
