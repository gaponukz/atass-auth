package controller

import (
	"auth/src/entities"
	"encoding/json"
	"net/http"
)

type credentials struct {
	Gmail    string `json:"gmail"`
	Password string `json:"password"`
}

type userCredentialsnDTO struct {
	FullName string `json:"fullName"`
	credentials
}

func decodeRequestBody(request *http.Request, data interface{}) error {
	err := json.NewDecoder(request.Body).Decode(data)

	if err != nil {
		return err
	}

	return nil
}

func getUserCredentialsFromBody(request *http.Request) (userCredentialsnDTO, error) {
	var creds userCredentialsnDTO

	err := decodeRequestBody(request, &creds)
	if err != nil {
		return userCredentialsnDTO{}, err
	}

	return creds, nil
}

func getGmailConfirmationFromBody(request *http.Request) (entities.FutureUser, error) {
	var creds entities.FutureUser

	err := decodeRequestBody(request, &creds)
	if err != nil {
		return entities.FutureUser{}, err
	}

	return creds, nil
}
