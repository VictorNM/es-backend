package user

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/victornm/es-backend/store"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	// Authentication error
	ErrNotAuthenticated = errors.New("not authenticated")

	// Registration error
	ErrEmailExisted = errors.New("email already existed")

	// Common error
	ErrInvalidInput = errors.New("invalid input")
)

type BasicSignInService interface {
	BasicSignIn(email, password string) (string, error)
}

type basicSignInService struct {
	userFinder Finder

	// jwt
	secret  string
	expired time.Duration
}

func NewBasicSignInService(userGetter Finder, secret string, expiredHours int) *basicSignInService {
	return &basicSignInService{
		userFinder: userGetter,
		secret:     secret,
		expired:    time.Duration(expiredHours) * time.Hour,
	}
}

// BasicSignIn use email and password for authentication
// Return a encrypted JWT token if sign in succeed
func (s *basicSignInService) BasicSignIn(email, password string) (string, error) {
	input := struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}{
		Email:    email,
		Password: password,
	}

	if err := validate.Struct(input); err != nil {
		return "", ErrInvalidInput
	}

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
		AuthDTO: &AuthDTO{UserID: u.ID},
	})

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", ErrNotAuthenticated
	}

	return tokenString, nil
}

type JWTParserService interface {
	ParseToken(tokenString string) (*AuthDTO, error)
}

type jwtParserService struct {
	secret string
}

func (s *jwtParserService) ParseToken(tokenString string) (*AuthDTO, error) {
	var claims jwtClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(s.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrNotAuthenticated
	}

	return claims.AuthDTO, nil
}

type AuthDTO struct {
	UserID int `json:"user_id"`
}

type jwtClaims struct {
	jwt.StandardClaims
	*AuthDTO
}

func NewJWTParserService(secret string) *jwtParserService {
	return &jwtParserService{secret: secret}
}

type RegisterService interface {
	Register(input *RegisterMutation) error
}

type RegisterMutation struct {
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,password"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type registerService struct {
	dao FindCreater
}

func (s *registerService) Register(input *RegisterMutation) error {
	if err := validate.Struct(input); err != nil {
		return ErrInvalidInput
	}

	if _, err := s.dao.FindUserByEmail(input.Email); err == nil {
		return ErrEmailExisted
	}

	hashed, err := hashPassword(input.Password)
	if err != nil {
		return err
	}

	_, err = s.dao.CreateUser(&store.UserRow{
		Email:          input.Email,
		HashedPassword: hashed,
		IsActive:       false,
	})

	// TODO: send email invitation

	return err
}

func NewRegisterService(finder FindCreater) *registerService {
	return &registerService{dao: finder}
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}
