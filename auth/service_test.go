package auth_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/victornm/es-backend/auth"
	"github.com/victornm/es-backend/store"
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
)

type mockUserDAO struct {
	users []*store.UserRow
}

func (dao *mockUserDAO) FindUserByEmail(email string) (*store.UserRow, error) {
	for _, u := range dao.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func genPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}

	return string(hashed)
}

func newMockUserDao(users []*store.UserRow) *mockUserDAO {
	return &mockUserDAO{users: users}
}

var users = []*store.UserRow{
	{
		ID:             1,
		Email:          "victornm@es.com",
		HashedPassword: genPassword("123abc"),
	},
}

func TestBasicLogin(t *testing.T) {
	tests := map[string]struct {
		Email       string
		Password    string
		WantedError error
	}{
		"happy login": {
			"victornm@es.com",
			"123abc",
			nil,
		},
		"email not existed": {
			"foo@bar.com",
			"123abc",
			auth.ErrNotAuthenticated,
		},
		"password not match": {
			"victornm@es.com",
			"xyz321",
			auth.ErrNotAuthenticated,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := auth.NewService(newMockUserDao(users), "#12345", 24)
			_, err := s.BasicSignIn(test.Email, test.Password)
			assert.Equal(t, test.WantedError, err)
		})
	}
}

func TestParseToken(t *testing.T) {
	t.Run("receive valid token", func(t *testing.T) {
		u := users[0]

		s := auth.NewService(newMockUserDao(users), "#12345", 24)
		tokenString, err := s.BasicSignIn(u.Email, "123abc")
		if err != nil {
			t.FailNow()
		}

		userAuth, err := s.ParseToken(tokenString)
		assert.NoError(t, err)
		assert.Equal(t, u.ID, userAuth.UserID)
	})
}
