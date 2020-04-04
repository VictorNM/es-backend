package user

import (
	"github.com/stretchr/testify/assert"
	"github.com/victornm/es-backend/store"
	"testing"
)

func TestBasicLogin(t *testing.T) {
	userInDB := []*store.UserRow{
		{
			Email:          "victornm@es.com",
			HashedPassword: mustHashPassword("1234abcd"),
		},
	}

	tests := map[string]struct {
		Email       string
		Password    string
		WantedError error
	}{
		"happy login": {
			"victornm@es.com",
			"1234abcd",
			nil,
		},
		"email not existed": {
			"foo@bar.com",
			"1234abcd",
			ErrNotAuthenticated,
		},
		"password not match": {
			"victornm@es.com",
			"4321bcda",
			ErrNotAuthenticated,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			dao := newMockUserDao()
			dao.seed(userInDB)

			s := NewBasicSignInService(dao, "#12345", 24)
			_, err := s.BasicSignIn(test.Email, test.Password)
			assert.Equal(t, test.WantedError, err)
		})
	}
}

func TestParseToken(t *testing.T) {
	usersInDB := []*store.UserRow{
		{
			Email:          "victornm@es.com",
			HashedPassword: mustHashPassword("1234abcd"),
		},
	}

	t.Run("receive valid token", func(t *testing.T) {
		dao := newMockUserDao()
		dao.seed(usersInDB)

		u := dao.users[0]
		secret := "#12345"

		s := NewBasicSignInService(dao, secret, 24)
		tokenString, err := s.BasicSignIn(u.Email, "1234abcd")
		if err != nil {
			t.FailNow()
		}

		parser := NewJWTParserService(secret)

		userAuth, err := parser.ParseToken(tokenString)
		assert.NoError(t, err)
		assert.Equal(t, u.ID, userAuth.UserID)
	})
}

func TestRegister(t *testing.T) {
	usersInDB := []*store.UserRow{
		{
			Email:          "victornm@es.com",
			HashedPassword: mustHashPassword("1234abcd"),
		},
	}

	tests := map[string]struct {
		input *RegisterMutation

		wantedErr error
	}{
		"happy case": {
			input: &RegisterMutation{
				Email:                "newEmail@gmail.com",
				Password:             "1234abcd",
				PasswordConfirmation: "1234abcd",
			},
			wantedErr: nil,
		},

		"existed email": {
			input: &RegisterMutation{
				Email:                "victornm@es.com",
				Password:             "1234abcd",
				PasswordConfirmation: "1234abcd",
			},
			wantedErr: ErrEmailExisted,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			dao := newMockUserDao()
			dao.seed(usersInDB)
			s := NewRegisterService(dao)

			// when
			err := s.Register(test.input)

			// then
			assert.Equal(t, test.wantedErr, err)
		})
	}
}

func TestRegister_ValidateInput(t *testing.T) {
	tests := map[string]*RegisterMutation{
		"invalid email": {
			Email:                "not an email",
			Password:             "abcd1234",
			PasswordConfirmation: "abcd1234",
		},

		"password length < 8": {
			Email:                "foo@bar.com",
			Password:             "123456",
			PasswordConfirmation: "123456",
		},

		"password length > 32": {
			Email:                "foo@bar.com",
			Password:             "abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
			PasswordConfirmation: "abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
		},

		"password not contain any letter": {
			Email:                "foo@bar.com",
			Password:             "12345678",
			PasswordConfirmation: "12345678",
		},

		"password not contain any digit": {
			Email:                "foo@bar.com",
			Password:             "abcdqwer",
			PasswordConfirmation: "abcdqwer",
		},

		"password confirmation not match": {
			Email:                "foo@bar.com",
			Password:             "abcd1234",
			PasswordConfirmation: "4321zyxw",
		},
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			dao := newMockUserDao()

			s := NewRegisterService(dao)

			err := s.Register(input)

			assert.Equal(t, ErrInvalidInput, err)
		})
	}
}
