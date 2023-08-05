package session

import (
	"auth/src/dto"
	"auth/src/errors"
	"auth/src/services/session"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	service := session.NewSessionService("test")
	payload := dto.CreateTokenDTO{RememberHim: true, UserInfoDTO: dto.UserInfoDTO{ID: "123"}}

	token, expir, err := service.CreateToken(payload)
	if err != nil {
		t.Error(err.Error())
	}

	if time.Now().After(expir) {
		t.Error("Token is already expired")
	}

	if token == "" {
		t.Error("Token is empty")
	}

	if len(token) < 13 {
		t.Error("token too short")
	}
}

func TestGetInfo(t *testing.T) {
	service := session.NewSessionService("test")
	payload := dto.CreateTokenDTO{RememberHim: true, UserInfoDTO: dto.UserInfoDTO{ID: "123"}}

	token, _, err := service.CreateToken(payload)
	if err != nil {
		t.Error(err.Error())
	}

	info, err := service.GetInfoFromToken(token)
	if err != nil {
		t.Error(err.Error())
	}

	if payload.AllowsAdvertisement != info.AllowsAdvertisement {
		t.Errorf("expected %t, got %t", payload.AllowsAdvertisement, info.AllowsAdvertisement)
	}

	if payload.ID != info.ID {
		t.Errorf("expected %s, got %s", payload.ID, info.ID)
	}
}

func TestGetInfoWrongWithToken(t *testing.T) {
	service := session.NewSessionService("test")
	payload := dto.CreateTokenDTO{
		RememberHim: true,
		UserInfoDTO: dto.UserInfoDTO{
			ID:    "123",
			Gmail: "test@example.com",
		},
	}

	token, _, err := service.CreateToken(payload)
	if err != nil {
		t.Error(err.Error())
	}

	info, err := service.GetInfoFromToken(token + "blabla")
	if err != nil {
		if err != errors.ErrUserNotFound {
			t.Error(err.Error())
		}
	} else {
		t.Error("Get info from wromg token without error")
	}

	if info.Gmail != "" {
		t.Errorf("find gmail: %s", info.Gmail)
	}
}

func TestUpdateToken(t *testing.T) {
	service := session.NewSessionService("test")
	payload := dto.CreateTokenDTO{RememberHim: true, UserInfoDTO: dto.UserInfoDTO{Phone: "123"}}

	token, _, err := service.CreateToken(payload)
	if err != nil {
		t.Error(err.Error())
	}

	newToken, expir, err := service.UpdateToken(token, dto.UpdateUserDTO{Phone: "321"})
	if err != nil {
		t.Error(err.Error())
	}

	if newToken == token {
		t.Error("token should be updated")
	}

	if time.Now().After(expir) {
		t.Error("New token is already expired")
	}

	info, err := service.GetInfoFromToken(newToken)
	if err != nil {
		t.Error(err.Error())
	}

	if "321" != info.Phone {
		t.Errorf("expected %s, got %s", "321", info.Phone)
	}
}
