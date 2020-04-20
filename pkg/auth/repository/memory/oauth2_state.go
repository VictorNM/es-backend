package memory

import (
	"errors"
	"github.com/victornm/es-backend/pkg/auth/internal"
)

// var GlobalOAuth2StateRepository = NewOauth2StateRepository()

type OAuth2StateRepository struct {
	data map[string]*internal.OAuth2State
}

func (r *OAuth2StateRepository) CreateState(state *internal.OAuth2State) error {
	if _, ok := r.data[state.Nonce]; ok {
		return errors.New("state already existed")
	}

	r.data[state.Nonce] = state
	return nil
}

func (r *OAuth2StateRepository) FindState(nonce string) (*internal.OAuth2State, error) {
	if state, ok := r.data[nonce]; ok {
		return state, nil
	}

	return nil, errors.New("state not found")
}

func NewOauth2StateRepository() *OAuth2StateRepository {
	return &OAuth2StateRepository{
		data: map[string]*internal.OAuth2State{},
	}
}
