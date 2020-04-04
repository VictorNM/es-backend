package user

import (
	"errors"
)

var ErrNotFound = errors.New("user not found")

type queryService struct {
	finder Finder
}

/*
 * USER PROFILE
 */
type GetProfileQuery interface {
	GetProfile(id int) (*ProfileDTO, error)
}

type ProfileDTO struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

func (s *queryService) GetProfile(id int) (*ProfileDTO, error) {
	u, err := s.finder.FindUserByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	return &ProfileDTO{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
	}, nil
}

func NewQueryService(finder Finder) *queryService {
	return &queryService{
		finder: finder,
	}
}
