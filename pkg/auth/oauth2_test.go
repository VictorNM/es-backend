package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/victornm/es-backend/pkg/auth"
	"github.com/victornm/es-backend/pkg/auth/mock"
)

func TestOAuth2Register(t *testing.T) {
	providerName := mock.ProviderName
	usersFromProvider := map[string]*User{
		"code_1": {Email: "foo@bar.com", Provider: providerName},
		"code_2": {Email: "admin@es.com", Provider: providerName},
	}

	usersInDB := []*User{
		{
			Email:          "admin@es.com",
			Username:       "admin",
			HashedPassword: MustHashPassword("1234abcd"),
		},
	}

	tests := map[string]struct {
		code string

		wantedErr error
	}{
		"happy": {
			code: "code_1",

			wantedErr: nil,
		},

		"email existed": {
			code: "code_2",

			wantedErr: ErrEmailExisted,
		},

		"provider return an error": {
			code: "some invalid code",

			wantedErr: ErrNotAuthenticated,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			provider := mock.NewOAuth2Provider()
			provider.Seed(usersFromProvider)

			repository := newUserRepository()
			repository.Seed(usersInDB)

			s := NewOAuth2Service(&OAuth2Config{
				UserRepository: repository,
				Providers:      []OAuth2Provider{provider},
			})

			err := s.OAuth2Register(OAuth2Input{
				Provider: providerName,
				Code:     test.code,
			})
			assertIsError(t, test.wantedErr, err)

			if test.wantedErr == nil {
				u, err := repository.FindUserByEmail(usersFromProvider[test.code].Email)
				assert.NoError(t, err)
				assert.Equal(t, providerName, u.Provider)
			}
		})
	}
}

func TestOAuth2SignIn(t *testing.T) {
	providerName := mock.ProviderName
	usersFromProvider := map[string]*User{
		"code_1": {Email: "admin@es.com", Provider: providerName},
		"code_2": {Email: "email_not_existed@bar.com", Provider: providerName},
		"code_3": {Email: "not_activated@es.com", Provider: providerName},
		"code_4": {Email: "provider_not_match@es.com", Provider: providerName},
	}

	usersInDB := []*User{
		{
			Email:          "admin@es.com",
			Username:       "admin",
			HashedPassword: MustHashPassword("1234abcd"),
			IsActive:       true,
			Provider:       providerName,
		},

		{
			Email:          "not_activated@es.com",
			Username:       "not_activated",
			HashedPassword: MustHashPassword("1234abcd"),
			IsActive:       false,
			Provider:       providerName,
		},

		{
			Email:          "provider_not_match@es.com",
			Username:       "provider_not_match",
			HashedPassword: MustHashPassword("1234abcd"),
			IsActive:       false,
			Provider:       "anotherProvider",
		},
	}

	tests := map[string]struct {
		code string

		wantedErr error
	}{
		"happy": {
			code: "code_1",

			wantedErr: nil,
		},

		"email not exist": {
			code: "code_2",

			wantedErr: ErrNotAuthenticated,
		},

		"user not activated": {
			code: "code_3",

			wantedErr: ErrNotActivated,
		},

		"provider not match": {
			code: "code_4",

			wantedErr: ErrNotAuthenticated,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			provider := mock.NewOAuth2Provider()
			provider.Seed(usersFromProvider)

			repository := newUserRepository()
			repository.Seed(usersInDB)

			s := NewOAuth2Service(&OAuth2Config{
				UserRepository: repository,
				Providers:      []OAuth2Provider{provider},
				JWTService:     NewJWTService("1234", 1),
			})

			token, err := s.OAuth2SignIn(OAuth2Input{
				Provider: providerName,
				Code:     test.code,
			})
			assertIsError(t, test.wantedErr, err)

			if test.wantedErr == nil {
				assert.NotEmpty(t, token)
			}
		})
	}
}
