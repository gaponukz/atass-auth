package controller

import (
	"auth/src/settings"
	"auth/src/storage"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Credentials struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
}

type Claims struct {
	FullName string `json:"fullName"`
	jwt.RegisteredClaims
}

type Controller struct {
	Storage  storage.IUserStorage
	Settings settings.Settings
}

func (contr *Controller) Singin(responseWriter http.ResponseWriter, request *http.Request) {
	var creds Credentials

	err := json.NewDecoder(request.Body).Decode(&creds)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedUser, err := contr.Storage.GetByGmail(creds.Gmail)

	if err != nil || expectedUser.Password != creds.Password {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		FullName: expectedUser.FullName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(contr.Settings.JwtSecret)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(responseWriter, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

}

func (contr *Controller) Refresh(responseWriter http.ResponseWriter, request *http.Request) {}

func (contr *Controller) Logout(responseWriter http.ResponseWriter, request *http.Request) {}

func (contr *Controller) Welcome(responseWriter http.ResponseWriter, request *http.Request) {}
