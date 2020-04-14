package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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
	readRepository ReadRepository

	// jwt
	secret  string
	expired time.Duration
}

func NewBasicSignInService(repository ReadRepository, secret string, expiredHours int) *basicSignInService {
	return &basicSignInService{
		readRepository: repository,
		secret:         secret,
		expired:        time.Duration(expiredHours) * time.Hour,
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

	u, err := s.readRepository.FindUserByEmail(email)
	if err != nil {
		return "", wrapError(ErrNotAuthenticated, "email")
	}

	err = u.ComparePassword(password)
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
		UserAuthDTO: &UserAuthDTO{UserID: u.ID},
	})

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", wrapError(ErrUnknown, err)
	}

	return tokenString, nil
}

// ===== JWT =====

type JWTParserService interface {
	ParseToken(tokenString string) (*UserAuthDTO, error)
}

type jwtParserService struct {
	secret string
}

func (s *jwtParserService) ParseToken(tokenString string) (*UserAuthDTO, error) {
	var claims jwtClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(s.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, wrapError(ErrNotAuthenticated, err)
	}

	return claims.UserAuthDTO, nil
}

type UserAuthDTO struct {
	UserID int `json:"user_id"`
}

type jwtClaims struct {
	jwt.StandardClaims
	*UserAuthDTO
}

func NewJWTParserService(secret string) *jwtParserService {
	return &jwtParserService{secret: secret}
}

// ===== Register =====

type RegisterService interface {
	Register(input *RegisterInput) error
}

type RegisterInput struct {
	Email                string `json:"email" validate:"required,email"`
	Username             string `json:"username" validate:"required,min=2,max=30"`
	Password             string `json:"password" validate:"required,password"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	FullName             string `json:"full_name" validate:"required"`
}

type registerService struct {
	repository Repository
	sender     ActivationEmailSender
}

func (s *registerService) Register(input *RegisterInput) error {
	if err := validate(input); err != nil {
		return wrapError(ErrInvalidInput, err)
	}

	if _, err := s.repository.FindUserByEmail(input.Email); err == nil {
		return wrapError(ErrEmailExisted, input.Email)
	}

	if _, err := s.repository.FindUserByUsername(input.Username); err == nil {
		return wrapError(ErrUsernameExisted, input.Username)
	}

	u, err := NewUser(input.Email, input.Password)
	if err != nil {
		return wrapError(ErrUnknown, err)
	}

	id, err := s.repository.CreateUser(u)

	if err != nil {
		return wrapError(ErrUnknown, err)
	}

	time.AfterFunc(time.Millisecond, func() {
		s.sender.SendActivationEmail(id)
	})

	return nil
}

type ActivationEmailSender interface {
	SendActivationEmail(userID int)
}

func NewRegisterService(repository Repository, sender ActivationEmailSender) *registerService {
	return &registerService{
		repository: repository,
		sender:     sender,
	}
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