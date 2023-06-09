package controller

import (
	"auth/src/entities"
	"auth/src/registration"
	"auth/src/resetPassword"
	"auth/src/settings"
	"auth/src/storage"
	"fmt"
	"net/http"
	"time"
)

type IUserStorage interface {
	Create(entities.User) error
	Delete(entities.User) error
	GetByGmail(string) (entities.User, error)
	UpdatePassword(entities.User, string) error
}

type Controller struct {
	Storage              IUserStorage
	Settings             settings.Settings
	RegistrationService  registration.RegistrationService
	ResetPasswordService resetPassword.ResetPasswordService
}

func (contr *Controller) Signin(responseWriter http.ResponseWriter, request *http.Request) {
	expectedUser, err := getRegisteredUserFromRequestBody(request, contr.Storage)

	if err != nil {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, expirationTime, err := getTemporaryToken(
		userInfoDTO{
			FullName:          expectedUser.FullName,
			RememberHim:       expectedUser.RememberHim,
			PurchasedRouteIds: expectedUser.PurchasedRouteIds,
		},
		contr.Settings.JwtSecret,
	)

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

	dto, err := getAuthorizedUserDataFromCookie(tokenCookie, contr.Settings.JwtSecret)

	if err != nil {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	response := fmt.Sprintf("Hi %s, remember? - %t", dto.FullName, dto.RememberHim)
	responseWriter.Write([]byte(response))
}

func (contr *Controller) Signup(responseWriter http.ResponseWriter, request *http.Request) {
	creds, err := getUserCredentialsFromBody(request)

	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	key, err := contr.RegistrationService.GetInformatedFutureUser(creds.Gmail)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = contr.RegistrationService.AddUserToFutureStorage(entities.FutureUser{
		UniqueKey: key,
		User: entities.User{
			Gmail:       creds.Gmail,
			Password:    creds.Password,
			FullName:    creds.FullName,
			Phone:       creds.Phone,
			RememberHim: creds.RememberHim,
		},
	})

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (contr *Controller) ConfirmRegistration(responseWriter http.ResponseWriter, request *http.Request) {
	user, err := getGmailConfirmationFromBody(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	err = contr.RegistrationService.RemoveUserFromFutureStorage(user)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	token, expirationTime, err := getTemporaryToken(
		userInfoDTO{
			FullName:          user.FullName,
			RememberHim:       user.RememberHim,
			PurchasedRouteIds: user.PurchasedRouteIds,
		},
		contr.Settings.JwtSecret,
	)

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

func (contr *Controller) Refresh(responseWriter http.ResponseWriter, request *http.Request) {
	tokenCookie, err := request.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	claims, tokenErr := parseClaimsFromToken(tokenCookie.Value, contr.Settings.JwtSecret)

	if tokenErr != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	newToken, expirationTime, newTokernErr := genarateNewTemporaryTokenFromClaims(claims, contr.Settings.JwtSecret)

	if newTokernErr != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(responseWriter, &http.Cookie{
		Name:    "token",
		Value:   newToken,
		Expires: expirationTime,
	})
}

func (contr *Controller) Logout(responseWriter http.ResponseWriter, request *http.Request) {
	http.SetCookie(responseWriter, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
}

func (contr *Controller) ResetPassword(responseWriter http.ResponseWriter, request *http.Request) {
	gmail, err := getGmailFromBody(request)

	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	code, err := contr.ResetPasswordService.GenerateAndSendCodeToGmail(gmail)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = contr.ResetPasswordService.AddUserToTemporaryStorage(storage.UserCredentials{
		Gmail: gmail,
		Key:   code,
	})
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (contr *Controller) ConfirmResetPassword(responseWriter http.ResponseWriter, request *http.Request) {
	user, err := getResetPasswordConfirmationFromBody(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	err = contr.ResetPasswordService.ChangeUserPassword(
		storage.UserCredentials{
			Gmail: user.Gmail,
			Key:   user.Key,
		},
		user.Password,
	)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	http.Redirect(responseWriter, request, "/signin_page", http.StatusFound)
}
