package auth_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/victornm/es-backend/auth"
	"github.com/victornm/es-backend/user"
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
)

type mockUserDAO struct {
	users []*user.DTO
}

func (dao *mockUserDAO) FindUserByEmail(email string) (*user.DTO, error) {
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

func newMockUserDao(users []*user.DTO) *mockUserDAO {
	return &mockUserDAO{users: users}
}

func TestBasicLogin(t *testing.T) {
	users := []*user.DTO{
		{
			Email:          "vinhnm@sendo.vn",
			HashedPassword: genPassword("123abc"),
		},
	}

	tests := map[string]struct {
		Email       string
		Password    string
		WantedError error
	}{
		"happy login": {
			"vinhnm@sendo.vn",
			"123abc",
			nil,
		},
		"email not existed": {
			"foo@bar.com",
			"123abc",
			auth.ErrNotAuthenticated,
		},
		"password not match": {
			"vinhnm@sendo.vn",
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
