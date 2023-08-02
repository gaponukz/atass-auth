package controller

import (
	"auth/src/dto"
	"auth/src/entities"
	"auth/src/errors"
	"net/http"
	"time"
)

type signinService interface {
	Login(string, string) (entities.UserEntity, error)
}

type signupService interface {
	SendGeneratedCode(string) (string, error)
	AddUserToTemporaryStorage(dto.GmailWithKeyPairDTO) error
	RegisterUserOnRightCode(dto.SignUpDTO) (string, error)
}

type settingsService interface {
	UpdateWithFields(id string, fields dto.UpdateUserDTO) error
}

type showUserRoutesService interface {
	ShowRoutes(id string) ([]entities.Path, error)
}

type resetPasswordService interface {
	NotifyUser(string) (string, error)
	AddUserToTemporaryStorage(dto.GmailWithKeyPairDTO) error
	CancelPasswordResetting(dto.GmailWithKeyPairDTO) error
	ChangeUserPassword(dto.PasswordResetDTO) error
}

type Controller struct {
	jwtSecret             string
	signinService         signinService
	signupService         signupService
	resetPasswordService  resetPasswordService
	settingsService       settingsService
	showUserRoutesService showUserRoutesService
}

func NewController(jwtSecret string, signinService signinService, signupService signupService,
	resetPasswordService resetPasswordService, settingsService settingsService, showUserRoutesService showUserRoutesService) *Controller {

	return &Controller{
		jwtSecret:             jwtSecret,
		signinService:         signinService,
		signupService:         signupService,
		resetPasswordService:  resetPasswordService,
		settingsService:       settingsService,
		showUserRoutesService: showUserRoutesService,
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
		if err == errors.ErrUserNotFound {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}

		responseWriter.WriteHeader(http.StatusInternalServerError)
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

	err = c.signupService.AddUserToTemporaryStorage(dto.GmailWithKeyPairDTO{
		Gmail: gmail,
		Key:   key,
	})
	if err != nil {
		if err == errors.ErrUserAlreadyExists {
			responseWriter.WriteHeader(http.StatusConflict)
			return
		}

		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c Controller) ConfirmRegistration(responseWriter http.ResponseWriter, request *http.Request) {
	newUser, err := getSignUpDto(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := c.signupService.RegisterUserOnRightCode(newUser)
	if err != nil {
		if err == errors.ErrRegisterRequestMissing {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, expirationTime, err := genarateToken(
		createTokenDTO{
			userInfoDTO: userInfoDTO{
				ID:                  id,
				Gmail:               newUser.Gmail,
				Phone:               newUser.Phone,
				FullName:            newUser.FullName,
				AllowsAdvertisement: newUser.AllowsAdvertisement,
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
	claims, tokenErr := getClaimsFromRequest(request, c.jwtSecret)
	if tokenErr != nil {
		responseWriter.WriteHeader(http.StatusUnauthorized)
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

	err = c.resetPasswordService.AddUserToTemporaryStorage(dto.GmailWithKeyPairDTO{
		Gmail: gmail,
		Key:   code,
	})
	if err != nil {
		if err == errors.ErrUserNotFound {
			responseWriter.WriteHeader(http.StatusNotFound)
			return
		}
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c Controller) CancelPasswordResetting(responseWriter http.ResponseWriter, request *http.Request) {
	pair, err := getGmailWithKeyDTO(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.resetPasswordService.CancelPasswordResetting(pair)
	if err != nil {
		if err == errors.ErrPasswordResetRequestMissing {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

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

	err = c.resetPasswordService.ChangeUserPassword(user)
	if err != nil {
		if err == errors.ErrPasswordResetRequestMissing {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

		responseWriter.WriteHeader(http.StatusInternalServerError)
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

func (c Controller) UpdateUserInfo(responseWriter http.ResponseWriter, request *http.Request) {
	dto, err := getUpdateUserDTO(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	oldClaims, tokenErr := getClaimsFromRequest(request, c.jwtSecret)
	if tokenErr != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.settingsService.UpdateWithFields(oldClaims.userInfoDTO.ID, dto)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}

	newClaims := &claims{
		userInfoDTO: userInfoDTO{
			ID:                  oldClaims.userInfoDTO.ID,
			Gmail:               oldClaims.userInfoDTO.Gmail,
			Phone:               dto.Phone,
			FullName:            dto.FullName,
			AllowsAdvertisement: dto.AllowsAdvertisement,
		},
		RegisteredClaims: oldClaims.RegisteredClaims,
	}

	newToken, expirationTime, newTokernErr := genarateTokenFromClaims(newClaims, c.jwtSecret)

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

func (c Controller) ShowUserRoutes(responseWriter http.ResponseWriter, request *http.Request) {
	userInfo, status := userInfoFromRequest(request, c.jwtSecret)
	if status != http.StatusOK {
		responseWriter.WriteHeader(int(status))
		return
	}

	routes, err := c.showUserRoutesService.ShowRoutes(userInfo.ID)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBytes, err := dumpsJson(routes)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonBytes)
}
