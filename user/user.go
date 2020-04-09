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
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	YearOfBirth int    `json:"year_of_birth"`
	Country     string `json:"country"`
	Gender      string `json:"gender"`
	Language    string `json:"language"`
}

func (s *queryService) GetProfile(id int) (*ProfileDTO, error) {
	u, err := s.finder.FindUserByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	return &ProfileDTO{
		ID:          u.ID,
		Email:       u.Email,
		FullName:    u.FullName,
		Phone:       u.Phone,
		YearOfBirth: u.YearOfBirth,
		Country:     u.Country,
		Gender:      u.Gender,
		Language:    u.Language,
	}, nil
}

func NewQueryService(finder Finder) *queryService {
	return &queryService{
		finder: finder,
	}
}
