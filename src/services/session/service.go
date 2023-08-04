package session

import (
	"auth/src/errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type userInfoDTO struct {
	ID                  string `json:"id"`
	Gmail               string `json:"gmail"`
	FullName            string `json:"fullName"`
	Phone               string `json:"phone"`
	AllowsAdvertisement bool   `json:"allowsAdvertisement"`
}

type createTokenDTO struct {
	RememberHim bool `json:"rememberHim"`
	userInfoDTO
}

type claims struct {
	userInfoDTO
	jwt.RegisteredClaims
}

type userSession struct {
	secret string
}

func (s userSession) CreateToken(data createTokenDTO) (string, time.Time, error) {
	expirationTime := s.getTokenExpirationTime(data.RememberHim)

	claims := &claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		userInfoDTO: data.userInfoDTO,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func (s userSession) GetInfoFromToken(token string) (userInfoDTO, error) {
	claims := &claims{}

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return userInfoDTO{}, fmt.Errorf("failed ParseWithClaims: %v", err)
	}

	if !jwtToken.Valid {
		return userInfoDTO{}, errors.ErrUserNotFound
	}

	return claims.userInfoDTO, nil
}

func (s userSession) UpdateToken(token string, info userInfoDTO) (string, time.Time, error) {
	_claims := &claims{}
	tkn, err := jwt.ParseWithClaims(token, _claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return "", time.Time{}, fmt.Errorf("Can not ParseWithClaims: %v", err)
	}

	if !tkn.Valid {
		return "", time.Time{}, errors.ErrUserNotFound
	}

	if time.Until(_claims.ExpiresAt.Time) > 30*time.Second {
		return token, _claims.ExpiresAt.Time, nil
	}

	_claims = &claims{
		userInfoDTO: userInfoDTO{
			ID:                  _claims.userInfoDTO.ID,
			Gmail:               _claims.userInfoDTO.Gmail,
			Phone:               info.Phone,
			FullName:            info.FullName,
			AllowsAdvertisement: info.AllowsAdvertisement,
		},
		RegisteredClaims: _claims.RegisteredClaims,
	}

	_claims.ExpiresAt = jwt.NewNumericDate(_claims.ExpiresAt.Time)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, _claims)

	newTokenString, err := newToken.SignedString([]byte(s.secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("Can not SignedString: %v", err)
	}

	return newTokenString, _claims.ExpiresAt.Time, nil
}

func (s userSession) getTokenExpirationTime(remember bool) time.Time {
	if remember {
		return time.Now().Add(24 * 11 * time.Hour)
	}

	return time.Now().Add(10 * time.Minute)
}
