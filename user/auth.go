package user

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/victornm/es-backend/store"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	// Authentication error
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrNotActivated     = errors.New("not activated")

	// Registration error
	ErrEmailExisted    = errors.New("email already existed")
	ErrUsernameExisted = errors.New("username already existed")

	// Common error
	ErrInvalidInput = errors.New("invalid input")
	ErrUnknown      = errors.New("unknown error")
)

// ===== Basic sign in =====

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

	if err := validate(input); err != nil {
		return "", wrapError(ErrInvalidInput, err)
	}

	u, err := s.userFinder.FindUserByEmail(email)
	if err != nil {
		return "", wrapError(ErrNotAuthenticated, "email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return "", wrapError(ErrNotAuthenticated, err)
	}

	if err != nil {
		return "", wrapError(ErrNotAuthenticated, err)
	}

	if !u.IsActive {
		return "", wrapError(ErrNotActivated)
	}

	// sign successfully
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
		return "", wrapError(ErrUnknown, err)
	}

	return tokenString, nil
}

// ===== JWT =====

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
		return nil, wrapError(ErrNotAuthenticated, err)
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

// ===== Register =====

type RegisterService interface {
	Register(input *RegisterMutation) error
}

type RegisterMutation struct {
	Email                string `json:"email" validate:"required,email"`
	Username             string `json:"username" validate:"required,min=2,max=30"`
	Password             string `json:"password" validate:"required,password"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	FullName             string `json:"full_name" validate:"required"`
}

type registerService struct {
	dao       FindCreator
	publisher Publisher
}

type Publisher interface {
	Publish(e interface{})
}

type Registered struct {
	UserID int
}

func (s *registerService) Register(input *RegisterMutation) error {
	if err := validate(input); err != nil {
		return wrapError(ErrInvalidInput, err)
	}

	if _, err := s.dao.FindUserByEmail(input.Email); err == nil {
		return wrapError(ErrEmailExisted, input.Email)
	}

	if _, err := s.dao.FindUserByUsername(input.Username); err == nil {
		return wrapError(ErrUsernameExisted, input.Username)
	}

	hashed, err := hashPassword(input.Password)
	if err != nil {
		return wrapError(ErrUnknown, err)
	}

	u := &store.UserRow{
		Email:          input.Email,
		HashedPassword: hashed,
		IsActive:       false,
		ActivationKey:  uuid.New().String(),
	}
	id, err := s.dao.CreateUser(u)

	if err != nil {
		return wrapError(ErrUnknown, err)
	}

	// TODO: send email invitation
	go s.publisher.Publish(Registered{UserID: id})

	return nil
}

func NewRegisterService(dao FindCreator, publisher Publisher) *registerService {
	return &registerService{
		dao:       dao,
		publisher: publisher,
	}
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func wrapError(err error, msgAndArgs ...interface{}) error {
	return fmt.Errorf("%w: %s", err, messageFromMsgAndArgs(msgAndArgs))
}

func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}
