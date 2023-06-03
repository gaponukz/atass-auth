package controller

import (
	"auth/src/settings"
	"auth/src/storage"
	"fmt"
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

func (contr *Controller) Welcome(responseWriter http.ResponseWriter, request *http.Request) {
	tokenCookie, err := request.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	userFullName, err := getAuthorizedUserDataFromCookie(tokenCookie, contr.Settings.JwtSecret)

	if err != nil {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	responseWriter.Write([]byte(fmt.Sprintf("Welcome %s!", userFullName)))
}

func (contr *Controller) Refresh(responseWriter http.ResponseWriter, request *http.Request) {}

func (contr *Controller) Logout(responseWriter http.ResponseWriter, request *http.Request) {}
