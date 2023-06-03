package controller

import (
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

func getUserCredentialsFromBody(request *http.Request) (userCredentialsnDTO, error) {
	var creds userCredentialsnDTO

	err := json.NewDecoder(request.Body).Decode(&creds)

	if err != nil {
		return userCredentialsnDTO{}, err
	}

	return creds, nil
}
