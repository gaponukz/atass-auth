package controller

import (
	"auth/src/application/dto"
	"auth/src/domain/entities"
	"auth/src/domain/errors"
	"net/http"
	"time"
)

type signinService interface {
	Login(string, string) (entities.UserEntity, error)
}

type sessionService interface {
	CreateToken(dto.CreateTokenDTO) (string, time.Time, error)
	GetInfoFromToken(string) (dto.UserInfoDTO, error)
	UpdateToken(string, dto.UpdateUserDTO) (string, time.Time, error)
	RefreshToken(string) (string, time.Time, error)
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
	signinService        signinService
	signupService        signupService
	resetPasswordService resetPasswordService
	sessionService       sessionService

	settingsService settingsService

	showUserRoutesService showUserRoutesService
}

func NewController(
	signinService signinService,
	signupService signupService,
	resetPasswordService resetPasswordService,
	settingsService settingsService,
	showUserRoutesService showUserRoutesService,
	sessionService sessionService,
) *Controller {

	return &Controller{
		signinService:         signinService,
		signupService:         signupService,
		resetPasswordService:  resetPasswordService,
		settingsService:       settingsService,
		showUserRoutesService: showUserRoutesService,
		sessionService:        sessionService,
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

	token, expirationTime, err := c.sessionService.CreateToken(
		dto.CreateTokenDTO{
			RememberHim: creds.RememberHim,
			UserInfoDTO: dto.UserInfoDTO{
				ID:                  user.ID,
				Gmail:               user.Gmail,
				FullName:            user.FullName,
				Phone:               user.Phone,
				AllowsAdvertisement: user.AllowsAdvertisement,
			},
		},
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
		if err == errors.ErrUserAlreadyExists {
			responseWriter.WriteHeader(http.StatusConflict)
			return
		}

		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = c.signupService.AddUserToTemporaryStorage(dto.GmailWithKeyPairDTO{
		Gmail: gmail,
		Key:   key,
	})
	if err != nil {
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

	token, expirationTime, err := c.sessionService.CreateToken(
		dto.CreateTokenDTO{
			UserInfoDTO: dto.UserInfoDTO{
				ID:                  id,
				Gmail:               newUser.Gmail,
				Phone:               newUser.Phone,
				FullName:            newUser.FullName,
				AllowsAdvertisement: newUser.AllowsAdvertisement,
			},
		},
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
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	newToken, expirationTime, err := c.sessionService.RefreshToken(tokenCookie.Value)
	if err != nil {
		if err == errors.ErrUserNotFound {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err == errors.ErrTokenEarlyToUpdate {
			responseWriter.WriteHeader(http.StatusBadRequest)
			return
		}

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
		if err == errors.ErrUserNotFound {
			responseWriter.WriteHeader(http.StatusNotFound)
			return
		}
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = c.resetPasswordService.AddUserToTemporaryStorage(dto.GmailWithKeyPairDTO{
		Gmail: gmail,
		Key:   code,
	})
	if err != nil {
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
	tokenCookie, err := request.Cookie("token")
	if err != nil {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	info, err := c.sessionService.GetInfoFromToken(tokenCookie.Value)
	if err != nil {
		if err == errors.ErrUserNotFound {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}

		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBytes, err := dumpsJson(info)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonBytes)
}

func (c Controller) UpdateUserInfo(responseWriter http.ResponseWriter, request *http.Request) {
	tokenCookie, err := request.Cookie("token")
	if err != nil {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	infoToUpdate, err := getUpdateUserDTO(request)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	newToken, expirationTime, err := c.sessionService.UpdateToken(tokenCookie.Value, infoToUpdate)
	if err != nil {
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
	tokenCookie, err := request.Cookie("token")
	if err != nil {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	info, err := c.sessionService.GetInfoFromToken(tokenCookie.Value)
	if err != nil {
		if err == errors.ErrUserNotFound {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}

		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	routes, err := c.showUserRoutesService.ShowRoutes(info.ID)
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
