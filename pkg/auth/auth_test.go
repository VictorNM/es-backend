package auth_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
	. "github.com/victornm/es-backend/pkg/auth"
	"github.com/victornm/es-backend/pkg/auth/mock"
	gatewayMemory "github.com/victornm/es-backend/pkg/store/memory"
)

func newUserRepository() *mock.AuthUserRepository {
	return mock.NewRepository(gatewayMemory.NewUserGateway())
}

func TestBasicSignIn(t *testing.T) {
	userInDB := []*User{
		{
			Email:          "victornm@es.com",
			Username:       "victornm@es.com",
			HashedPassword: MustHashPassword("1234abcd"),
			IsActive:       true,
		},
		{
			Email:          "nguyenmauvinh@es.com",
			Username:       "nguyenmauvinh@es.com",
			HashedPassword: MustHashPassword("1234abcd"),
			IsActive:       false,
		},
	}

	tests := map[string]struct {
		email     string
		password  string
		wantedErr error
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
		"user not activated": {
			"nguyenmauvinh@es.com",
			"1234abcd",
			ErrNotActivated,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			repository := newUserRepository()
			repository.Seed(userInDB)

			s := New(&Config{
				UserRepository: repository,
				JWTService:     NewJWTService("#12345", 24),
			})

			_, err := s.BasicSignIn(test.email, test.password)
			assertIsError(t, test.wantedErr, err)
		})
	}
}

func TestRegister(t *testing.T) {
	usersInDB := []*User{
		{
			Email:          "victornm@es.com",
			Username:       "victorNM",
			HashedPassword: MustHashPassword("1234abcd"),
		},
	}

	tests := map[string]struct {
		input *RegisterInput

		wantedErr error
	}{
		"happy case": {
			input: &RegisterInput{
				Email:                "newEmail@gmail.com",
				Username:             "newUser",
				Password:             "1234abcd",
				PasswordConfirmation: "1234abcd",
				FullName:             "VictorNM",
			},
			wantedErr: nil,
		},

		"existed email": {
			input: &RegisterInput{
				Email:                "victornm@es.com",
				Username:             "newUser",
				Password:             "1234abcd",
				PasswordConfirmation: "1234abcd",
				FullName:             "VictorNM",
			},
			wantedErr: ErrEmailExisted,
		},

		"existed username": {
			input: &RegisterInput{
				Email:                "newEmail@gmail.com",
				Username:             "victornm",
				Password:             "1234abcd",
				PasswordConfirmation: "1234abcd",
				FullName:             "VictorNM",
			},
			wantedErr: ErrUsernameExisted,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			repository := newUserRepository()
			repository.Seed(usersInDB)

			s := New(&Config{
				UserRepository: repository,
				Mailer:         &mock.Mailer{},
				ActivateURL:    "/activate",
			})

			// when
			err := s.Register(test.input)

			// then
			assertIsError(t, test.wantedErr, err)
		})
	}
}

func TestRegister_SendActivationMail(t *testing.T) {
	// given
	repository := newUserRepository()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	sent := false
	mailer := &mock.Mailer{SendFunc: func(subject string, tmpl string, data interface{}, to []string) error {
		sent = true
		wg.Done()
		return nil
	}}

	s := New(&Config{
		UserRepository: repository,
		Mailer:         mailer,
		ActivateURL:    "/activation",
	})

	// when
	err := s.Register(&RegisterInput{
		Email:                "newEmail@gmail.com",
		Username:             "newUser",
		Password:             "1234abcd",
		PasswordConfirmation: "1234abcd",
		FullName:             "VictorNM",
	})
	wg.Wait()

	// then
	require.NoError(t, err)
	assert.True(t, sent)
}

func TestRegister_ValidateInput(t *testing.T) {
	t.Run("valid inputs", func(t *testing.T) {
		tests := []*RegisterInput{
			{
				Email:                "foo@bar.com",
				Username:             "lucifer_silver",
				Password:             "abcd1234",
				PasswordConfirmation: "abcd1234",
				FullName:             "Nguyen Mau Vinh",
			},
		}

		for i, input := range tests {
			t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
				dao := gatewayMemory.NewUserGateway()
				repository := mock.NewRepository(dao)

				s := New(&Config{
					UserRepository: repository,
					Mailer:         &mock.Mailer{},
				})

				assert.NoError(t, s.Register(input))
			})
		}
	})

	t.Run("invalid inputs", func(t *testing.T) {
		tests := map[string]*RegisterInput{
			"invalid email": {
				Email:                "not an email",
				Username:             "lucifer_silver",
				Password:             "abcd1234",
				PasswordConfirmation: "abcd1234",
				FullName:             "Nguyen Mau Vinh",
			},

			"password length < 8": {
				Email:                "foo@bar.com",
				Username:             "lucifer_silver",
				Password:             "123456",
				PasswordConfirmation: "123456",
				FullName:             "Nguyen Mau Vinh",
			},

			"password not contain any letter": {
				Email:                "foo@bar.com",
				Username:             "lucifer_silver",
				Password:             "12345678",
				PasswordConfirmation: "12345678",
				FullName:             "Nguyen Mau Vinh",
			},

			"password not contain any digit": {
				Email:                "foo@bar.com",
				Username:             "lucifer_silver",
				Password:             "abcdqwer",
				PasswordConfirmation: "abcdqwer",
				FullName:             "Nguyen Mau Vinh",
			},

			"password confirmation not match": {
				Email:                "foo@bar.com",
				Username:             "lucifer_silver",
				Password:             "abcd1234",
				PasswordConfirmation: "4321zyxw",
				FullName:             "Nguyen Mau Vinh",
			},

			"username < 2": {
				Email:                "not an email",
				Username:             "s",
				Password:             "abcd1234",
				PasswordConfirmation: "abcd1234",
				FullName:             "Nguyen Mau Vinh",
			},

			"empty full name": {
				Email:                "not an email",
				Username:             "lucifer_silver",
				Password:             "abcd1234",
				PasswordConfirmation: "abcd1234",
				FullName:             "",
			},
		}

		for name, input := range tests {
			t.Run(name, func(t *testing.T) {
				err := input.Valid()

				assert.Error(t, err)
			})
		}
	})
}

func assertIsError(t *testing.T, wanted, got error) {
	t.Helper()
	if !errors.Is(got, wanted) {
		t.Errorf("Error %v is not an %v", got, wanted)
	}
}
