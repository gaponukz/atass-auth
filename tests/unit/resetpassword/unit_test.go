package resetpassword

import (
	"auth/src/application/dto"
	"auth/src/application/usecases/passreset"
	"auth/src/domain/entities"
	"auth/src/domain/errors"
	"auth/src/utils"
	"auth/tests/unit/mocks"
	"testing"
)

func TestNotifyUser(t *testing.T) {
	const expectedCode = "12345"

	sm := mocks.NewMockStorage()
	notify := mocks.NewMockUserNotifier()
	sendCode := func() string { return expectedCode }
	hash := func(s string) string { return s }

	s := passreset.NewResetPasswordService(sm, nil, notify, hash, sendCode)

	_, err := s.NotifyUser("user@gmail.com")
	if err != errors.ErrUserNotFound {
		t.Error("Can notify not registered gmail")
	}

	testUser := dto.CreateUserDTO{Gmail: "user@gmail.com"}

	_, err = sm.Create(testUser)
	if err != nil {
		t.Fatal(err)
	}

	code, err := s.NotifyUser("user@gmail.com")
	if err != nil {
		t.Error(err.Error())
	}
	if code != expectedCode {
		t.Errorf("Expected code %s, got %s", expectedCode, code)
	}
}

func TestAddUserToTemporaryStorage(t *testing.T) {
	sm := mocks.NewMockStorage()
	tsm := mocks.NewTemporaryStorageMock()
	s := passreset.NewResetPasswordService(sm, tsm, nil, nil, nil)

	testUser := dto.CreateUserDTO{Gmail: "user@gmail.com"}
	pair := dto.GmailWithKeyPairDTO{Gmail: "user@gmail.com", Key: "12345"}

	_, err := sm.Create(testUser)
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

func TestCancelPasswordResetting(t *testing.T) {
	const gmail = "test@gmain.com"
	const key = "12345"
	testPair := dto.GmailWithKeyPairDTO{Gmail: gmail, Key: key}
	testUser := dto.CreateUserDTO{Gmail: gmail, Password: "old"}

	sm := mocks.NewMockStorage()
	tsm := mocks.NewTemporaryStorageMock()
	s := passreset.NewResetPasswordService(sm, tsm, nil, nil, nil)

	err := tsm.Create(testPair)
	if err != nil {
		t.Fatal(err)
	}

	_, err = sm.Create(testUser)
	if err != nil {
		t.Fatal(err)
	}

	err = s.CancelPasswordResetting(testPair)
	if err != nil {
		t.Errorf("error changing password: %v", err)
	}

	_, err = tsm.GetByUniqueKey(key)
	if err == nil {
		t.Error("after cancel password reseting pair still in temp storage")
	}
}

func TestChangeUserPasswordWithWrongCode(t *testing.T) {
	const gmail = "test@gmain.com"
	const key = "12345"
	testPair := dto.GmailWithKeyPairDTO{Gmail: gmail, Key: key}
	testUser := dto.CreateUserDTO{Gmail: gmail, Password: "old"}

	sm := mocks.NewMockStorage()
	tsm := mocks.NewTemporaryStorageMock()
	hash := func(s string) string { return s }
	s := passreset.NewResetPasswordService(sm, tsm, nil, hash, nil)

	err := tsm.Create(testPair)
	if err != nil {
		t.Fatal(err)
	}

	_, err = sm.Create(testUser)
	if err != nil {
		t.Fatal(err)
	}

	err = s.ChangeUserPassword(dto.PasswordResetDTO{Gmail: gmail, Key: key + "lala", Password: "new"})
	if err == nil {
		t.Error("can change password with wrong code")
	}
}

func TestChangeUserPassword(t *testing.T) {
	const gmail = "test@gmain.com"
	const key = "12345"
	testPair := dto.GmailWithKeyPairDTO{Gmail: gmail, Key: key}
	testUser := dto.CreateUserDTO{Gmail: gmail, Password: "old"}

	sm := mocks.NewMockStorage()
	tsm := mocks.NewTemporaryStorageMock()
	hash := func(s string) string { return s }
	s := passreset.NewResetPasswordService(sm, tsm, nil, hash, nil)

	err := tsm.Create(testPair)
	if err != nil {
		t.Fatal(err)
	}

	_, err = sm.Create(testUser)
	if err != nil {
		t.Fatal(err)
	}

	err = s.ChangeUserPassword(dto.PasswordResetDTO{Gmail: gmail, Key: key, Password: "new"})
	if err != nil {
		t.Errorf("error changing password: %v", err)
	}

	_, err = tsm.GetByUniqueKey(key)
	if err == nil {
		t.Error("after password reseting pair still in temp storage")
	}

	users, err := sm.ReadAll()
	if err != nil {
		t.Error(err.Error())
	}

	user, err := utils.Find(users, func(user entities.User) bool {
		return user.Gmail == gmail
	})
	if err != nil {
		t.Errorf("error getting user after reseting: %v", err)
	}

	if user.Password != hash("new") {
		t.Error("password not correct after changing")
	}
}
