package controller

import (
	"auth/src/entities"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type claims struct {
	FullName string `json:"fullName"`
	jwt.RegisteredClaims
}

type GetByGmailAbleStorage interface {
	GetByGmail(string) (entities.User, error)
}

func getRegisteredUserFromRequestBody(request *http.Request, storage GetByGmailAbleStorage) (entities.User, error) {
	creds, err := getUserCredentialsFromBody(request)

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

func getAuthorizedUserDataFromCookie(cookie *http.Cookie, jwtSecret string) (string, error) {
	claims := &claims{}

	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if !token.Valid {
		return "", fmt.Errorf("unauthorized")
	}

	return claims.FullName, err
}
