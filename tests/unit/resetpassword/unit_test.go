package resetpassword

import (
	"auth/src/entities"
	"auth/src/resetPassword"
	"auth/tests/unit/mocks"
	"testing"
)

func TestNotifyUser(t *testing.T) {
	const expectedCode = "12345"

	s := resetPassword.NewResetPasswordService(nil, nil, func(gmail string, key string) error {
		return nil
	}, func() string {
		return expectedCode
	})

	code, err := s.NotifyUser("user@gmail.com")
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
	s := resetPassword.NewResetPasswordService(sm, tsm, nil, nil)

	testUser := entities.User{Gmail: "user@gmail.com"}
	pair := entities.GmailWithKeyPair{Gmail: "user@gmail.com", Key: "12345"}

	err := sm.Create(testUser)
	if err != nil {
		t.Fatal(err)
	}

	err = s.AddUserToTemporaryStorage(pair)
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

func TestChangeUserPassword(t *testing.T) {
	const gmail = "test@gmain.com"
	const key = "12345"
	testPair := entities.GmailWithKeyPair{Gmail: gmail, Key: key}
	testUser := entities.User{Gmail: gmail, Password: "old"}

	sm := mocks.NewMockStorage()
	tsm := mocks.NewTemporaryStorageMock()
	s := resetPassword.NewResetPasswordService(sm, tsm, nil, nil)

	err := tsm.Create(testPair)
	if err != nil {
		t.Fatal(err)
	}

	err = sm.Create(testUser)
	if err != nil {
		t.Fatal(err)
	}

	err = s.ChangeUserPassword(testPair, "new")
	if err != nil {
		t.Errorf("error changing password: %v", err)
	}

	_, err = tsm.GetByUniqueKey(key)
	if err == nil {
		t.Error("after password reseting pair still in temp storage")
	}

	user, err := sm.GetByGmail(gmail)
	if err != nil {
		t.Errorf("error getting user after reseting: %v", err)
	}

	if user.Password != "new" {
		t.Errorf("password not changed: expected new, got %s", user.Password)
	}
}
