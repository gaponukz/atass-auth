package controller

import (
	"auth/src/application/dto"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func decodeRequestBody(request *http.Request, data interface{}) error {
	err := json.NewDecoder(request.Body).Decode(data)

	if err != nil {
		return err
	}

	return nil
}

func getGmailWithKeyDTO(request *http.Request) (dto.GmailWithKeyPairDTO, error) {
	var data dto.GmailWithKeyPairDTO

	err := decodeRequestBody(request, &data)
	return data, err
}

func getSignInDto(request *http.Request) (dto.SignInDTO, error) {
	var creds dto.SignInDTO

	err := decodeRequestBody(request, &creds)
	return creds, err
}

func getSignUpDto(request *http.Request) (dto.SignUpDTO, error) {
	var creds dto.SignUpDTO

	err := decodeRequestBody(request, &creds)
	return creds, err
}

func getPasswordResetDto(request *http.Request) (dto.PasswordResetDTO, error) {
	var creds dto.PasswordResetDTO

	err := decodeRequestBody(request, &creds)
	return creds, err
}

func getUpdateUserDTO(request *http.Request) (dto.UpdateUserDTO, error) {
	var dto dto.UpdateUserDTO

	err := decodeRequestBody(request, &dto)
	return dto, err
}

func getOneStringFieldFromBody(request *http.Request, field string) (string, error) {
	value, err := getOneFieldFromBody(request, field)
	stringValue, ok := value.(string)
	if !ok && err == nil {
		return "", fmt.Errorf("expected type string, got %T", value)
	}

	return stringValue, err
}

func getOneFieldFromBody(request *http.Request, field string) (interface{}, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	value, exists := data[field]
	if !exists {
		return nil, fmt.Errorf("%s field not found", field)
	}

	return value, nil
}

func dumpsJson(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return []byte(""), err
	}

	return jsonData, nil
}
