package controller

import (
	"auth/src/entities"
	"net/http"
	"time"
)

type signinService interface {
	Login(string, string) (entities.UserEntity, error)
	UserProfile(string) (entities.UserEntity, error)
}

type signupService interface {
	SendGeneratedCode(string) (string, error)
	AddUserToTemporaryStorage(entities.GmailWithKeyPair) error
	RegisterUserOnRightCode(entities.GmailWithKeyPair, entities.User) (string, error)
}

type settingsService interface {
	UpdateWithFields(string, interface{}) error
	SubscribeUserToRoutes(string, string) error
}

type resetPasswordService interface {
	NotifyUser(string) (string, error)
	AddUserToTemporaryStorage(entities.GmailWithKeyPair) error
	CancelPasswordResetting(entities.GmailWithKeyPair) error
	ChangeUserPassword(entities.GmailWithKeyPair, string) error
}

type Controller struct {
	jwtSecret            string
	signinService        signinService
	signupService        signupService
	resetPasswordService resetPasswordService
	settingsService      settingsService
}

func NewController(jwtSecret string, signinService signinService, signupService signupService,
	resetPasswordService resetPasswordService, settingsService settingsService) *Controller {

	return &Controller{
		jwtSecret:            jwtSecret,
		signinService:        signinService,
		signupService:        signupService,
		resetPasswordService: resetPasswordService,
		settingsService:      settingsService,
	}
}

func (c Controller) Signin(responseWriter http.ResponseWriter, request *http.Request) {
	creds, err := getSignInDto(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := c.signinService.Login(creds.Gmail, creds.Password)
	if err != nil {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, expirationTime, err := genarateToken(
		createTokenDTO{
			RememberHim: creds.RememberHim,
			userInfoDTO: userInfoDTO{
				ID:                  user.ID,
				Gmail:               user.Gmail,
				FullName:            user.FullName,
				Phone:               user.Phone,
				AllowsAdvertisement: user.AllowsAdvertisement,
			},
		},
		c.jwtSecret,
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
	gmail, err := getOneStringFieldFromBody(request, "gmail")

	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	key, err := c.signupService.SendGeneratedCode(gmail)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = c.signupService.AddUserToTemporaryStorage(entities.GmailWithKeyPair{
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

	id, err := c.signupService.RegisterUserOnRightCode(entities.GmailWithKeyPair{
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
			userInfoDTO: userInfoDTO{
				ID:                  id,
				Gmail:               dto.Gmail,
				Phone:               dto.Phone,
				FullName:            dto.FullName,
				AllowsAdvertisement: dto.AllowsAdvertisement,
			},
		},
		c.jwtSecret,
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

	claims, tokenErr := getClaimsFromToken(tokenCookie.Value, c.jwtSecret)
	if tokenErr != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	newToken, expirationTime, newTokernErr := genarateTokenFromClaims(claims, c.jwtSecret)

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
	gmail, err := getOneStringFieldFromBody(request, "gmail")

	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	code, err := c.resetPasswordService.NotifyUser(gmail)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = c.resetPasswordService.AddUserToTemporaryStorage(entities.GmailWithKeyPair{
		Gmail: gmail,
		Key:   code,
	})
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c Controller) CancelPasswordResetting(responseWriter http.ResponseWriter, request *http.Request) {
	user, err := getPasswordResetDto(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.resetPasswordService.CancelPasswordResetting(
		entities.GmailWithKeyPair{
			Gmail: user.Gmail,
			Key:   user.Key,
		},
	)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (c Controller) ConfirmResetPassword(responseWriter http.ResponseWriter, request *http.Request) {
	user, err := getPasswordResetDto(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.resetPasswordService.ChangeUserPassword(
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

func (c Controller) GetUserInfo(responseWriter http.ResponseWriter, request *http.Request) {
	userInfo, status := userInfoFromRequest(request, c.jwtSecret)
	if status != http.StatusOK {
		responseWriter.WriteHeader(int(status))
		return
	}

	jsonBytes, err := dumpsJson(userInfo)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonBytes)
}

func (c Controller) SubscribeToTheRoute(responseWriter http.ResponseWriter, request *http.Request) {
	routeId, err := getOneStringFieldFromBody(request, "routeId")
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	userInfo, status := userInfoFromRequest(request, c.jwtSecret)
	if status != http.StatusOK {
		responseWriter.WriteHeader(int(status))
		return
	}

	err = c.settingsService.SubscribeUserToRoutes(userInfo.ID, routeId)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c Controller) UpdateUserInfo(responseWriter http.ResponseWriter, request *http.Request) {
	dto, err := getUpdateUserDTO(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	user, status := userInfoFromRequest(request, c.jwtSecret)
	if status != http.StatusOK {
		responseWriter.WriteHeader(int(status))
		return
	}

	err = c.settingsService.UpdateWithFields(user.ID, dto)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}
}
