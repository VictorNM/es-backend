package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/victornm/es-backend/user"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrNotAuthenticated = errors.New("not authenticated")
)

type BasicSignIner interface {
	BasicSignIn(email, password string) (string, error)
}

type UserFinder interface {
	FindUserByEmail(email string) (*user.DTO, error)
}

type service struct {
	userFinder UserFinder

	// jwt
	secret  string
	expired time.Duration
}

func NewService(userGetter UserFinder, secretKey string, expiredHour int) *service {
	return &service{
		userFinder: userGetter,
		secret:     secretKey,
		expired:    time.Duration(expiredHour) * time.Hour,
	}
}

// BasicSignIn use email and password for authentication
// Return a encrypted JWT token if sign in succeed
func (s *service) BasicSignIn(email, password string) (string, error) {
	u, err := s.userFinder.FindUserByEmail(email)
	if err != nil {
		return "", ErrNotAuthenticated
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return "", ErrNotAuthenticated
	}

	// unknown error
	if err != nil {
		return "", ErrNotAuthenticated
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(s.expired).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "auth.service",
		Subject:   u.Email,
	})

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", ErrNotAuthenticated
	}

	return tokenString, nil
}
