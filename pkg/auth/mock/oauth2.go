package mock

import (
	"errors"
	"github.com/victornm/es-backend/pkg/auth"
)

const ProviderName = "mock"

func NewOAuth2Provider() *mockProvider {
	return &mockProvider{users: map[string]*auth.User{}}
}

type mockProvider struct {
	users map[string]*auth.User
}

func (p *mockProvider) Name() string {
	return ProviderName
}

func (p *mockProvider) GetUser(code string) (*auth.User, error) {
	if u, ok := p.users[code]; ok {
		return u, nil
	}

	return nil, errors.New("invalid code")
}

func (p *mockProvider) Seed(users map[string]*auth.User) {
	for code, u := range users {
		p.users[code] = u
	}
}

func (p *mockProvider) Clear() {
	p.users = map[string]*auth.User{}
}
