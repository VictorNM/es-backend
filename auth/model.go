package auth

import (
	"github.com/google/uuid"
	"github.com/victornm/es-backend/store"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	*store.UserRow
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}

func NewUser(email, password string) (*User, error) {
	hashed, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		UserRow: &store.UserRow{
			Email:          email,
			HashedPassword: hashed,
			IsActive:       false,
			ActivationKey:  uuid.New().String(),
		},
	}, nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}