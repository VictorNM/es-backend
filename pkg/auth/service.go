package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/victornm/es-backend/pkg/auth/internal"
	"github.com/victornm/es-backend/pkg/errorutil"
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
	readRepository ReadUserRepository
	jwtService     JWTService
}

func NewBasicSignInService(repository ReadUserRepository, jwtService JWTService) *basicSignInService {
	return &basicSignInService{
		readRepository: repository,
		jwtService:     jwtService,
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

	if err := internal.Validate(input); err != nil {
		return "", errorutil.Wrap(ErrInvalidInput, err)
	}

	u, err := s.readRepository.FindUserByEmail(email)
	if err != nil {
		return "", errorutil.Wrap(ErrNotAuthenticated, "email")
	}

	err = u.ComparePassword(password)
	if err != nil {
		return "", errorutil.Wrap(ErrNotAuthenticated, err)
	}

	if !u.IsActive {
		return "", errorutil.Wrap(ErrNotActivated)
	}

	// sign successfully
	return s.jwtService.GenerateToken(u)
}

// ===== JWT =====

type JWTService interface {
	ParseToken(tokenString string) (*UserAuthDTO, error)
	GenerateToken(u *User) (string, error)
}

type jwtService struct {
	secret  string
	expired time.Duration
}

func (s *jwtService) GenerateToken(u *User) (string, error) {
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
		return "", errorutil.Wrap(ErrUnknown, err)
	}

	return tokenString, nil
}

func (s *jwtService) ParseToken(tokenString string) (*UserAuthDTO, error) {
	var claims jwtClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(s.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errorutil.Wrap(ErrNotAuthenticated, err)
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

func NewJWTService(secret string, expiredHours int) *jwtService {
	return &jwtService{
		secret:  secret,
		expired: time.Duration(expiredHours) * time.Hour,
	}
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
	repository UserRepository
	sender     ActivationEmailSender
}

func (s *registerService) Register(input *RegisterInput) error {
	if err := internal.Validate(input); err != nil {
		return errorutil.Wrap(ErrInvalidInput, err)
	}

	if _, err := s.repository.FindUserByEmail(input.Email); err == nil {
		return errorutil.Wrap(ErrEmailExisted, input.Email)
	}

	if _, err := s.repository.FindUserByUsername(input.Username); err == nil {
		return errorutil.Wrap(ErrUsernameExisted, input.Username)
	}

	u, err := NewUser(input.Email, input.Password)
	if err != nil {
		return errorutil.Wrap(ErrUnknown, err)
	}

	id, err := s.repository.CreateUser(u)

	if err != nil {
		return errorutil.Wrap(ErrUnknown, err)
	}

	time.AfterFunc(time.Millisecond, func() {
		s.sender.SendActivationEmail(id)
	})

	return nil
}

func NewRegisterService(repository UserRepository, sender ActivationEmailSender) *registerService {
	return &registerService{
		repository: repository,
		sender:     sender,
	}
}
