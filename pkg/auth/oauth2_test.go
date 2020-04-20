package auth_test

import (
	"github.com/stretchr/testify/assert"
	. "github.com/victornm/es-backend/pkg/auth"
	"github.com/victornm/es-backend/pkg/auth/repository/memory"
	memory2 "github.com/victornm/es-backend/pkg/store/memory"
	"net/url"
	"testing"
)

func TestOAuth2ConfigFactory(t *testing.T) {
	factory := NewOAuth2ClientFactory(WithGoogle("1234", "", ""))

	_, err := factory.GetService("facebook")
	assertIsError(t, ErrInvalidOAuth2Provider, err)

	_, err = factory.GetService("google")
	assert.NoError(t, err)
}

func TestOAuth2Register(t *testing.T) {
	factory := NewOAuth2ClientFactory(WithGoogle("1234", "", ""))

	stateRepository := memory.NewOauth2StateRepository()
	userRepository := memory.NewRepository(memory2.NewUserGateway())

	s := NewOAuth2RegisterService(stateRepository, userRepository, factory)

	authURL, err := s.OAuth2Register("google")
	assert.NoError(t, err)

	u, err := url.Parse(authURL)
	assert.NoError(t, err)
	assert.Equal(t, "1234", u.Query().Get("client_id"))
}