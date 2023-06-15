package controller

import (
	"auth/src/entities"
	"auth/src/registration"
	"auth/src/resetPassword"
	"auth/src/settings"
	"net/http"
	"time"
)

type userStorage interface {
	Create(entities.User) error
	Delete(entities.User) error
	GetByGmail(string) (entities.User, error)

	UpdatePassword(entities.User, string) error
	SubscribeToTheRoute(entities.User, string) error
}

type Controller struct {
	Storage              userStorage
	Settings             settings.Settings
	RegistrationService  registration.RegistrationService
	ResetPasswordService resetPassword.ResetPasswordService
}

func (contr *Controller) Signin(responseWriter http.ResponseWriter, request *http.Request) {
	creds, err := getSignInCredentialsFromBody(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	if !isUserAuthenticated(credentials{Gmail: creds.Gmail, Password: creds.Password}, contr.Storage) {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err != nil {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, expirationTime, err := getTemporaryToken(
		createTokenDTO{
			Gmail:       creds.Gmail,
			RememberHim: creds.RememberHim,
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

func (contr *Controller) Signup(responseWriter http.ResponseWriter, request *http.Request) {
	gmail, err := getGmailFromBody(request)

	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	key, err := contr.RegistrationService.GetInformatedFutureUser(gmail)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = contr.RegistrationService.AddUserToTemporaryStorage(entities.GmailWithKeyPair{
		Gmail: gmail,
		Key:   key,
	})

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (contr *Controller) ConfirmRegistration(responseWriter http.ResponseWriter, request *http.Request) {
	dto, err := getUserCredentialsFromBody(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	err = contr.RegistrationService.RegisterUserOnRightCode(entities.GmailWithKeyPair{
		Gmail: dto.Gmail,
		Key:   dto.Key,
	}, entities.User{
		Gmail:               dto.Gmail,
		Password:            dto.Password,
		Phone:               dto.Phone,
		FullName:            dto.FullName,
		AllowsAdvertisement: dto.AllowsAdvertisement,
	})
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	token, expirationTime, err := getTemporaryToken(
		createTokenDTO{
			Gmail: dto.Gmail,
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

	err = contr.ResetPasswordService.AddUserToTemporaryStorage(entities.GmailWithKeyPair{
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
		entities.GmailWithKeyPair{
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

func (contr *Controller) GetFullUserInfo(responseWriter http.ResponseWriter, request *http.Request) {
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

	fullUserInfo, err := contr.Storage.GetByGmail(dto.Gmail)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	fullUserInfo.Password = ""

	jsonBytes, err := loadStructIntoJson(fullUserInfo)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonBytes)
}

func (contr *Controller) SubscribeToTheRoute(responseWriter http.ResponseWriter, request *http.Request) {
	routeId, err := getRouteIdFromBody(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

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

	user, err := contr.Storage.GetByGmail(dto.Gmail)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = contr.Storage.SubscribeToTheRoute(user, routeId)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}
