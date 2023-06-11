package controller

import (
	"auth/src/entities"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type userInfoDTO struct {
	Gmail       string `json:"gmail"`
	RememberHim bool   `json:"rememberHim"`
}

type claims struct {
	userInfoDTO
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

func getTemporaryToken(infoDto userInfoDTO, jwtSecret string) (string, time.Time, error) {
	expirationTime := getExpirationTime(infoDto.RememberHim)

	claims := &claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		userInfoDTO: infoDto,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func getAuthorizedUserDataFromCookie(cookie *http.Cookie, jwtSecret string) (userInfoDTO, error) {
	claims := &claims{}

	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if !token.Valid {
		return userInfoDTO{}, fmt.Errorf("unauthorized")
	}

	dto := userInfoDTO{
		Gmail:       claims.Gmail,
		RememberHim: claims.RememberHim,
	}

	return dto, err
}

func parseClaimsFromToken(token, secret string) (*claims, error) {
	_claims := &claims{}
	tkn, err := jwt.ParseWithClaims(token, _claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return &claims{}, err
	}

	if !tkn.Valid {
		return &claims{}, fmt.Errorf("token is not valid")
	}

	return _claims, nil
}

func genarateNewTemporaryTokenFromClaims(oldClaims *claims, secret string) (string, time.Time, error) {
	expirationTime := getExpirationTime(oldClaims.RememberHim)
	oldClaims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, oldClaims)
	tockenStr, err := token.SignedString([]byte(secret))

	return tockenStr, expirationTime, err
}
