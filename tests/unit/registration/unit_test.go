package registration

import (
	"auth/src/dto"
	"auth/src/entities"
	"auth/src/services/signup"
	"auth/src/utils"
	"auth/tests/unit/mocks"
	"testing"
)

func TestSendGeneratedCode(t *testing.T) {
	const expectedCode = "12345"

	notify := func(gmail string, key string) error { return nil }
	generateCode := func() string { return expectedCode }
	hash := func(s string) string { return s }
	s := signup.NewRegistrationService(nil, nil, notify, generateCode, hash)

	code, err := s.SendGeneratedCode("user@gmail.com")
	if err != nil {
		t.Errorf("Error sending generated code: %v", err)
	}

	if code != expectedCode {
		t.Errorf("Expected code %s, got %s", expectedCode, code)
	}
}

func TestAddUserToTemporaryStorage(t *testing.T) {
	sm := mocks.NewMockStorage()
	tsm := mocks.NewTemporaryStorageMock()
	s := signup.NewRegistrationService(sm, tsm, nil, nil, nil)
	testUser := dto.GmailWithKeyPairDTO{Gmail: "user@gmail.com", Key: "12345"}

	err := s.AddUserToTemporaryStorage(testUser)
	if err != nil {
		t.Errorf("Error adding to temporary storage: %v", err)
	}

	user, err := tsm.GetByUniqueKey("12345")
	if err != nil {
		t.Errorf("after AddUserToTemporaryStorage user now adding: %v", err)
	}

	if user.Gmail != testUser.Gmail {
		t.Errorf("user in temporary expected %s, got %s", testUser.Gmail, user.Gmail)
	}
}

func TestAddAlreadyRegisteredUserToTemporaryStorage(t *testing.T) {
	sm := mocks.NewMockStorage()
	tsm := mocks.NewTemporaryStorageMock()
	s := signup.NewRegistrationService(sm, tsm, nil, nil, nil)
	testUser := dto.GmailWithKeyPairDTO{Gmail: "user@gmail.com", Key: "12345"}

	_, err := sm.Create(entities.User{Gmail: "user@gmail.com"})
	if err != nil {
		t.Fatal(err)
	}

	err = s.AddUserToTemporaryStorage(testUser)
	if err == nil {
		t.Error("successfully added already registered user")
	}
}

func TestAddWrongCodeUser(t *testing.T) {
	const expectedCode = "12345"
	sm := mocks.NewMockStorage()
	tsm := mocks.NewTemporaryStorageMock()
	generateCode := func() string { return expectedCode }
	s := signup.NewRegistrationService(sm, tsm, nil, generateCode, nil)
	testUser := dto.GmailWithKeyPairDTO{Gmail: "user@gmail.com", Key: expectedCode + "lala"}

	err := s.AddUserToTemporaryStorage(testUser)
	if err == nil {
		t.Error("successfully added registered user with wrong code")
	}
}

func TestRegisterUserOnRightCode(t *testing.T) {
	const expectedCode = "12345"
	sm := mocks.NewMockStorage()
	tsm := mocks.NewTemporaryStorageMock()
	generateCode := func() string { return expectedCode }
	hash := func(s string) string { return s }

	s := signup.NewRegistrationService(sm, tsm, nil, generateCode, hash)

	const testGmail = "test@gmail.com"
	pair := dto.GmailWithKeyPairDTO{Gmail: testGmail, Key: "12345"}
	user := entities.User{Gmail: testGmail}

	err := tsm.Create(pair)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.RegisterUserOnRightCode(dto.SignUpDTO{Gmail: testGmail, Key: "12345"})
	if err != nil {
		t.Errorf("RegisterUserOnRightCode error: %v", err)
	}

	_, err = tsm.GetByUniqueKey("12345")
	if err == nil {
		t.Error("pair still awiable in temporary storage")
	}

	users, err := sm.ReadAll()
	if err != nil {
		t.Error(err.Error())
	}

	u, err := utils.Find(users, func(user entities.UserEntity) bool {
		return user.Gmail == testGmail
	})
	if err != nil {
		t.Errorf("GetByGmail error: %v", err)
	}

	if u.Gmail != testGmail {
		t.Errorf("after register expected %s, got %s", testGmail, user.Gmail)
	}
}
