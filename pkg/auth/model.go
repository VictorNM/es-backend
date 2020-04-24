package auth

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID             int
	Email          string
	Username       string
	HashedPassword string
	FullName       string
	IsActive       bool
	IsSuperAdmin   bool
	ActivationKey  string
	Provider       string
	CreatedAt      time.Time
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}

// TODO: Add Username, FullName
func NewUser(email, password string) (*User, error) {
	hashed, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:          email,
		HashedPassword: hashed,
		IsActive:       false,
		IsSuperAdmin:   false,
		ActivationKey:  uuid.New().String(),
	}, nil
}

func NewOAuth2User(email, provider string) *User {
	return &User{
		Email:        email,
		IsActive:     true,
		IsSuperAdmin: false,
		Provider:     provider,
	}
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}