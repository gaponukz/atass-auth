package session

import (
	"auth/src/dto"
	"auth/src/errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type claims struct {
	dto.UserInfoDTO
	jwt.RegisteredClaims
}

type userSession struct {
	secret string
}

func NewUserSession(secret string) userSession {
	return userSession{secret: secret}
}

func (s userSession) CreateToken(data dto.CreateTokenDTO) (string, time.Time, error) {
	expirationTime := s.getTokenExpirationTime(data.RememberHim)

	claims := &claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		UserInfoDTO: data.UserInfoDTO,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func (s userSession) GetInfoFromToken(token string) (dto.UserInfoDTO, error) {
	claims := &claims{}

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return dto.UserInfoDTO{}, errors.ErrUserNotFound
	}

	if !jwtToken.Valid {
		return dto.UserInfoDTO{}, errors.ErrUserNotFound
	}

	return claims.UserInfoDTO, nil
}

func (s userSession) UpdateToken(token string, info dto.UpdateTokenDTO) (string, time.Time, error) {
	oldClaims, err := s.getClaimsFromToken(token)
	if err != nil {
		return "", time.Time{}, err
	}

	newClaims := &claims{
		UserInfoDTO: dto.UserInfoDTO{
			ID:                  oldClaims.UserInfoDTO.ID,
			Gmail:               oldClaims.UserInfoDTO.Gmail,
			Phone:               info.Phone,
			FullName:            info.FullName,
			AllowsAdvertisement: info.AllowsAdvertisement,
		},
		RegisteredClaims: oldClaims.RegisteredClaims,
	}

	return s.genarateTokenFromClaims(newClaims)
}

func (s userSession) RefreshToken(token string) (string, time.Time, error) {
	claims, err := s.getClaimsFromToken(token)
	if err != nil {
		return "", time.Time{}, errors.ErrUserNotFound
	}

	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return "", time.Time{}, errors.ErrTokenEarlyToUpdate
	}

	return s.genarateTokenFromClaims(claims)
}

func (s userSession) getTokenExpirationTime(remember bool) time.Time {
	if remember {
		return time.Now().Add(24 * 11 * time.Hour)
	}

	return time.Now().Add(10 * time.Minute)
}

func (s userSession) genarateTokenFromClaims(oldClaims *claims) (string, time.Time, error) {
	expirationTime := s.getTokenExpirationTime(false)
	oldClaims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, oldClaims)
	tockenStr, err := token.SignedString([]byte(s.secret))

	return tockenStr, expirationTime, err
}

func (s userSession) getClaimsFromToken(token string) (*claims, error) {
	_claims := &claims{}
	tkn, err := jwt.ParseWithClaims(token, _claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return &claims{}, err
	}

	if !tkn.Valid {
		return &claims{}, fmt.Errorf("token is not valid")
	}

	return _claims, nil
}
