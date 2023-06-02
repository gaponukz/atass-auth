package controller

import (
	"auth/src/settings"
	"auth/src/storage"
	"net/http"
)

type Controller struct {
	Storage  storage.IUserStorage
	Settings settings.Settings
}

func (contr *Controller) Singin(responseWriter http.ResponseWriter, request *http.Request) {
	expectedUser, err := getRegisteredUserFromRequestBody(request, contr.Storage)

	if err != nil {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, expirationTime, err := getTemporaryToken(expectedUser.FullName, contr.Settings.JwtSecret)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(responseWriter, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})

}

func (contr *Controller) Refresh(responseWriter http.ResponseWriter, request *http.Request) {}

func (contr *Controller) Logout(responseWriter http.ResponseWriter, request *http.Request) {}

func (contr *Controller) Welcome(responseWriter http.ResponseWriter, request *http.Request) {}
