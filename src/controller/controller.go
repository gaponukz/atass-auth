package controller

import (
	"auth/src/entities"
	"auth/src/settings"
	"net/http"
	"time"
)

type userStorage interface {
	Create(entities.User) (entities.UserEntity, error)
	ReadAll() ([]entities.UserEntity, error)
	ByID(string) (entities.UserEntity, error)
	Update(entities.UserEntity) error
	Delete(string) error
}

type registrationService interface {
	SendGeneratedCode(string) (string, error)
	AddUserToTemporaryStorage(entities.GmailWithKeyPair) error
	RegisterUserOnRightCode(entities.GmailWithKeyPair, entities.User) (string, error)
}

type resetPasswordService interface {
	NotifyUser(string) (string, error)
	AddUserToTemporaryStorage(entities.GmailWithKeyPair) error
	ChangeUserPassword(entities.GmailWithKeyPair, string) error
}

type Controller struct {
	Storage              userStorage
	Settings             settings.Settings
	RegistrationService  registrationService
	ResetPasswordService resetPasswordService
}

func (c Controller) Signin(responseWriter http.ResponseWriter, request *http.Request) {
	creds, err := getSignInDto(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := getIDifCredsValid(credentials{Gmail: creds.Gmail, Password: creds.Password}, c.Storage)
	if err != nil {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, expirationTime, err := genarateToken(
		createTokenDTO{
			ID:          id,
			RememberHim: creds.RememberHim,
		},
		c.Settings.JwtSecret,
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

func (c Controller) Signup(responseWriter http.ResponseWriter, request *http.Request) {
	gmail, err := getGmailFromBody(request)

	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	key, err := c.RegistrationService.SendGeneratedCode(gmail)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = c.RegistrationService.AddUserToTemporaryStorage(entities.GmailWithKeyPair{
		Gmail: gmail,
		Key:   key,
	})

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c Controller) ConfirmRegistration(responseWriter http.ResponseWriter, request *http.Request) {
	dto, err := getSignUpDto(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := c.RegistrationService.RegisterUserOnRightCode(entities.GmailWithKeyPair{
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

	token, expirationTime, err := genarateToken(
		createTokenDTO{
			ID: id,
		},
		c.Settings.JwtSecret,
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

func (c Controller) Refresh(responseWriter http.ResponseWriter, request *http.Request) {
	tokenCookie, err := request.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	claims, tokenErr := getClaimsFromToken(tokenCookie.Value, c.Settings.JwtSecret)

	if tokenErr != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	newToken, expirationTime, newTokernErr := genarateTokenFromClaims(claims, c.Settings.JwtSecret)

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

func (c Controller) Logout(responseWriter http.ResponseWriter, request *http.Request) {
	http.SetCookie(responseWriter, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
}

func (c Controller) ResetPassword(responseWriter http.ResponseWriter, request *http.Request) {
	gmail, err := getGmailFromBody(request)

	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	code, err := c.ResetPasswordService.NotifyUser(gmail)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = c.ResetPasswordService.AddUserToTemporaryStorage(entities.GmailWithKeyPair{
		Gmail: gmail,
		Key:   code,
	})
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c Controller) ConfirmResetPassword(responseWriter http.ResponseWriter, request *http.Request) {
	user, err := getPasswordResetDto(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.ResetPasswordService.ChangeUserPassword(
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
}

func (c Controller) GetFullUserInfo(responseWriter http.ResponseWriter, request *http.Request) {
	id, status := idFromRequest(request, c.Settings.JwtSecret)
	if status != http.StatusOK {
		responseWriter.WriteHeader(int(status))
		return
	}

	fullUserInfo, err := c.Storage.ByID(id)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	fullUserInfo.Password = ""

	jsonBytes, err := dumpsJson(fullUserInfo)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonBytes)
}

func (c Controller) SubscribeToTheRoute(responseWriter http.ResponseWriter, request *http.Request) {
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

	dto, err := getAuthorizedUserDataFromCookie(tokenCookie, c.Settings.JwtSecret)
	if err != nil {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := c.Storage.ByID(dto.ID)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.PurchasedRouteIds = append(user.PurchasedRouteIds, routeId)

	err = c.Storage.Update(user)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c Controller) ChangeUserName(responseWriter http.ResponseWriter, request *http.Request) {
	name, err := getOneStringFieldFromBody(request, "gmail")
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	id, status := idFromRequest(request, c.Settings.JwtSecret)
	if status != http.StatusOK {
		responseWriter.WriteHeader(int(status))
		return
	}

	user, err := c.Storage.ByID(id)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.FullName = name

	err = c.Storage.Update(user)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}
}
