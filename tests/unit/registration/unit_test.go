package registration

import (
	"auth/src/application/dto"
	"auth/src/application/usecases/signup"
	"auth/src/domain/entities"
	"auth/src/domain/errors"
	"auth/src/utils"
	"auth/tests/unit/mocks"
	"testing"
)

func TestPasswordValidation(t *testing.T) {
	s := signup.NewRegistrationService(nil, nil, nil, nil, nil)
	cases := []struct {
		password string
		isValid  bool
	}{
		{"qweasdzxcQ1", true},
		{"12345", false},
		{"Ab1$", false},
		{"aB1$cdEf", true},
		{"P@ss", false},
		{"ABCD1234", false},
		{"p@ssword", false},
		{"PASSWORD1234", false},
		{"P@ssw0rd1234567890", true},
		{"P@ssw0rd123", true},
		{"Speci1al$Character", true},
	}

	for _, test := range cases {
		if s.IsPasswordValid(test.password) != test.isValid {
			t.Errorf("For %s got %t expected %t", test.password, !test.isValid, test.isValid)
		}
	}
}

func TestPhoneValidations(t *testing.T) {
	s := signup.NewRegistrationService(nil, nil, nil, nil, nil)
	cases := []struct {
		password string
		isValid  bool
	}{
		{"o984516456", false},
		{"0984516456", true},
		{"380984516456", true},
		{"-380984516456", false},
		{"+380 66 538 29 59", true},
		{"5823946528346592384652375", false},
	}

	for _, test := range cases {
		if s.IsPhoneNumberValid(test.password) != test.isValid {
			t.Errorf("For %s got %t expected %t", test.password, !test.isValid, test.isValid)
		}
	}
}

func TestFullNameValidations(t *testing.T) {
	s := signup.NewRegistrationService(nil, nil, nil, nil, nil)
	cases := []struct {
		password string
		isValid  bool
	}{
		{"", false},
		{"Alex Was", true},
		{"Sam", false},
		{"Sam A", false},
		{"Sam Awrt", true},
		{"Sa Wa", true},
	}

	for _, test := range cases {
		if s.IsFullNameValid(test.password) != test.isValid {
			t.Errorf("For %s got %t expected %t", test.password, !test.isValid, test.isValid)
		}
	}
}

func TestSendGeneratedCode(t *testing.T) {
	const expectedCode = "12345"

	sm := mocks.NewMockStorage()
	notify := mocks.NewMockGmailNotifier()
	generateCode := func() string { return expectedCode }
	hash := func(s string) string { return s }
	s := signup.NewRegistrationService(sm, nil, notify, generateCode, hash)

	_, err := sm.Create(dto.CreateUserDTO{Gmail: "user@gmail.com"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.SendGeneratedCode("user@gmail.com")
	if err != errors.ErrUserAlreadyExists {
		t.Error("successfully added already registered user")
	}

	code, err := s.SendGeneratedCode("user2@gmail.com")
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

func TestAddWrongCodeUser(t *testing.T) {
	const expectedCode = "12345"
	sm := mocks.NewMockStorage()
	tsm := mocks.NewTemporaryStorageMock()
	generateCode := func() string { return expectedCode }
	hash := func(s string) string { return s }

	s := signup.NewRegistrationService(sm, tsm, nil, generateCode, hash)

	const testGmail = "test@gmail.com"
	pair := dto.GmailWithKeyPairDTO{Gmail: testGmail, Key: "12345"}

	err := tsm.Create(pair)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.RegisterUserOnRightCode(dto.SignUpDTO{Gmail: testGmail, Key: "wrongkey"})
	if err == nil {
		t.Error("can register with wrong code")
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

	err := tsm.Create(pair)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.RegisterUserOnRightCode(dto.SignUpDTO{Gmail: testGmail, Key: "12345", Password: "qweasdzxcQ1", FullName: "So Va", Phone: "0984516456"})
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

	u, err := utils.Find(users, func(user entities.User) bool {
		return user.Gmail == testGmail
	})
	if err != nil {
		t.Errorf("GetByGmail error: %v", err)
	}

	if u.Gmail != testGmail {
		t.Errorf("after register expected %s, got %s", testGmail, u.Gmail)
	}
}
