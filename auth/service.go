package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrNotAuthenticated = errors.New("not authenticated")
)

type BasicSignIner interface {
	BasicSignIn(email, password string) (string, error)
}

type TokenParser interface {
	ParseToken(tokenString string) (*UserAuth, error)
}

type service struct {
	userFinder UserFinder

	// jwt
	secret  string
	expired time.Duration
}

func NewService(userGetter UserFinder, secret string, expiredHours int) *service {
	return &service{
		userFinder: userGetter,
		secret:     secret,
		expired:    time.Duration(expiredHours) * time.Hour,
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.expired).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "auth.service",
		},
		UserAuth: &UserAuth{UserID: u.ID},
	})

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", ErrNotAuthenticated
	}

	return tokenString, nil
}

func (s *service) ParseToken(tokenString string) (*UserAuth, error) {
	var claims jwtClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(s.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrNotAuthenticated
	}

	return claims.UserAuth, nil
}

type UserAuth struct {
	UserID int `json:"user_id"`
}

type jwtClaims struct {
	jwt.StandardClaims
	*UserAuth
}
