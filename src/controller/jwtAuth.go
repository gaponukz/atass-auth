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

func getTemporaryToken(dto createTokenDTO, jwtSecret string) (string, time.Time, error) {
	expirationTime := getExpirationTime(dto.RememberHim)

	claims := &claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		userInfoDTO: userInfoDTO{Gmail: dto.Gmail},
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
		Gmail: claims.Gmail,
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
	expirationTime := getExpirationTime(false)
	oldClaims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, oldClaims)
	tockenStr, err := token.SignedString([]byte(secret))

	return tockenStr, expirationTime, err
}
