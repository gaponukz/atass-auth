package registration

import (
	"auth/src/entities"
	"auth/src/registration"
	"testing"
)

func TestSendGeneratedCode(t *testing.T) {
	const expectedCode = "12345"

	s := registration.NewRegistrationService(nil, nil, func(gmail string, key string) error {
		return nil
	}, func() string {
		return expectedCode
	})

	code, err := s.SendGeneratedCode("user@gmail.com")
	if err != nil {
		t.Errorf("Error sending generated code: %v", err)
	}

	if code != expectedCode {
		t.Errorf("Expected code %s, got %s", expectedCode, code)
	}
}

func TestAddUserToTemporaryStorage(t *testing.T) {
	sm := NewMockStorage()
	tsm := NewTemporaryStorageMock()
	s := registration.NewRegistrationService(sm, tsm, nil, nil)
	testUser := entities.GmailWithKeyPair{Gmail: "user@gmail.com", Key: "12345"}

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
	sm := NewMockStorage()
	tsm := NewTemporaryStorageMock()
	s := registration.NewRegistrationService(sm, tsm, nil, nil)
	testUser := entities.GmailWithKeyPair{Gmail: "user@gmail.com", Key: "12345"}

	err := sm.Create(entities.User{Gmail: "user@gmail.com"})
	if err != nil {
		t.Fatal(err)
	}

	err = s.AddUserToTemporaryStorage(testUser)
	if err == nil {
		t.Error("successfully added already registered user")
	}
}

func TestRegisterUserOnRightCode(t *testing.T) {
	sm := NewMockStorage()
	tsm := NewTemporaryStorageMock()
	s := registration.NewRegistrationService(sm, tsm, nil, func() string { return "12345" })

	const testGmail = "test@gmail.com"
	pair := entities.GmailWithKeyPair{Gmail: testGmail, Key: "12345"}
	user := entities.User{Gmail: testGmail}

	err := tsm.Create(pair)
	if err != nil {
		t.Fatal(err)
	}

	err = s.RegisterUserOnRightCode(pair, user)

	if err != nil {
		t.Errorf("RegisterUserOnRightCode error: %v", err)
	}

	_, err = tsm.GetByUniqueKey("12345")
	if err == nil {
		t.Error("pair still awiable in temporary storage")
	}

	u, err := sm.GetByGmail(testGmail)
	if err != nil {
		t.Errorf("GetByGmail error: %v", err)
	}

	if u.Gmail != testGmail {
		t.Errorf("after register expected %s, got %s", testGmail, user.Gmail)
	}

}
