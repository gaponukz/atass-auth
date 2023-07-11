package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type claims struct {
	userInfoDTO
	jwt.RegisteredClaims
}

func genarateToken(dto createTokenDTO, jwtSecret string) (string, time.Time, error) {
	expirationTime := getExpirationTime(dto.RememberHim)

	claims := &claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		userInfoDTO: dto.userInfoDTO,
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

	return claims.userInfoDTO, err
}

func getClaimsFromRequest(request *http.Request, secret string) (*claims, error) {
	tokenCookie, err := request.Cookie("token")

	if err != nil {
		return nil, err
	}

	_claims := &claims{}
	tkn, err := jwt.ParseWithClaims(tokenCookie.Value, _claims, func(token *jwt.Token) (interface{}, error) {
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

func genarateTokenFromClaims(oldClaims *claims, secret string) (string, time.Time, error) {
	expirationTime := getExpirationTime(false)
	oldClaims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, oldClaims)
	tockenStr, err := token.SignedString([]byte(secret))

	return tockenStr, expirationTime, err
}

type statusCode int

func userInfoFromRequest(request *http.Request, secret string) (userInfoDTO, statusCode) {
	tokenCookie, err := request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return userInfoDTO{}, http.StatusUnauthorized
		}

		return userInfoDTO{}, http.StatusBadRequest
	}

	dto, err := getAuthorizedUserDataFromCookie(tokenCookie, secret)
	if err != nil {
		return userInfoDTO{}, http.StatusUnauthorized
	}

	return dto, http.StatusOK
}
