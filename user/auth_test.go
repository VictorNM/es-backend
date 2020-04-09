package user

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/victornm/es-backend/event"
	"github.com/victornm/es-backend/store"
	"log"
	"sync"
	"testing"
)

func TestBasicLogin(t *testing.T) {
	userInDB := []*store.UserRow{
		{
			Email:          "victornm@es.com",
			HashedPassword: mustHashPassword("1234abcd"),
			IsActive:       true,
		},
		{
			Email:          "nguyenmauvinh@es.com",
			HashedPassword: mustHashPassword("1234abcd"),
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
			withoutValidate(func() {
				dao := newMockUserDao()
				dao.seed(userInDB)

				s := NewBasicSignInService(dao, "#12345", 24)
				_, err := s.BasicSignIn(test.email, test.password)
				assertIsError(t, test.wantedErr, err)

				log.Println(err)
			})
		})
	}
}

func TestParseToken(t *testing.T) {
	usersInDB := []*store.UserRow{
		{
			Email:          "victornm@es.com",
			HashedPassword: mustHashPassword("1234abcd"),
			IsActive:       true,
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

// ===== Register =====

func TestRegister(t *testing.T) {
	usersInDB := []*store.UserRow{
		{
			Email:          "victornm@es.com",
			Username:       "victorNM",
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
				Username:             "newUser",
				Password:             "1234abcd",
				PasswordConfirmation: "1234abcd",
			},
			wantedErr: nil,
		},

		"existed email": {
			input: &RegisterMutation{
				Email:                "victornm@es.com",
				Username:             "newUser",
				Password:             "1234abcd",
				PasswordConfirmation: "1234abcd",
			},
			wantedErr: ErrEmailExisted,
		},

		"existed username": {
			input: &RegisterMutation{
				Email:                "newEmail@gmail.com",
				Username:             "victornm",
				Password:             "1234abcd",
				PasswordConfirmation: "1234abcd",
			},
			wantedErr: ErrUsernameExisted,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			withoutValidate(func() {
				// given
				dao := newMockUserDao()
				dao.seed(usersInDB)
				s := NewRegisterService(dao, event.GetBus())

				// when
				err := s.Register(test.input)

				// then
				assertIsError(t, test.wantedErr, err)
			})
		})
	}
}

func TestRegister_PublishEvent(t *testing.T) {
	withoutValidate(func() {
		// given
		dao := newMockUserDao()
		bus := event.NewBus()
		s := NewRegisterService(dao, bus)
		var userRegisteredEvent Registered

		c := make(chan interface{})
		bus.Subscribe(Registered{}, c)
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go func() {
			for {
				select {
				case e := <-c:
					userRegisteredEvent = e.(Registered)
					wg.Done()
				}
			}
		}()

		// when
		err := s.Register(&RegisterMutation{
			Email:                "newEmail@gmail.com",
			Username:             "newUser",
			Password:             "1234abcd",
			PasswordConfirmation: "1234abcd",
		})
		wg.Wait()

		// Then
		assert.NoError(t, err)
		assert.NotEmpty(t, userRegisteredEvent.UserID)

		u, _ := dao.FindUserByID(userRegisteredEvent.UserID)
		assert.Equal(t, u.Email, "newEmail@gmail.com")
	})
}

func TestRegister_ValidateInput(t *testing.T) {
	t.Run("valid inputs", func(t *testing.T) {
		tests := []*RegisterMutation{
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
				dao := newMockUserDao()

				s := NewRegisterService(dao, event.GetBus())

				assert.NoError(t, s.Register(input))
			})
		}
	})

	t.Run("invalid inputs", func(t *testing.T) {
		tests := map[string]*RegisterMutation{
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
				dao := newMockUserDao()

				s := NewRegisterService(dao, event.GetBus())

				err := s.Register(input)

				assertIsError(t, ErrInvalidInput, err)
			})
		}
	})
}

func assertIsError(t *testing.T, wanted, got error) {
	if !errors.Is(got, wanted) {
		t.Errorf("Error %v is not an %v", got, wanted)
	}
}
