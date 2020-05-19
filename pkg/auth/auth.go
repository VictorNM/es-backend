package auth

import (
	"fmt"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"

	"github.com/victornm/es-backend/pkg/errorutil"
)

var _ Service = (*service)(nil)

type Service interface {
	BasicSignIn(email, password string) (string, error)
	Register(input *RegisterInput) error
}

type Config struct {
	UserRepository UserRepository

	JWTService JWTService

	ActivateURL string
	Mailer      Mailer
}

type service struct {
	userRepository UserRepository
	jwtService     JWTService
	sender         *activationEmailSender
}

func New(config *Config) Service {
	s := &service{
		userRepository: config.UserRepository,
		jwtService:     config.JWTService,
		sender:         &activationEmailSender{mailer: config.Mailer, repository: config.UserRepository, path: config.ActivateURL},
	}

	return s
}

// BasicSignIn use email and password for authentication
// Return a encrypted JWT token if sign in succeed
func (s *service) BasicSignIn(email, password string) (string, error) {
	input := &SignInInput{
		Email:    email,
		Password: password,
	}

	if err := validate(input); err != nil {
		return "", errorutil.Wrap(ErrInvalidInput, err)
	}

	u, err := s.userRepository.FindUserByEmail(email)
	if err != nil {
		return "", errorutil.Wrap(ErrNotAuthenticated, "email")
	}

	err = u.comparePassword(password)
	if err != nil {
		return "", errorutil.Wrap(ErrNotAuthenticated, err)
	}

	if !u.IsActive {
		return "", errorutil.Wrap(ErrNotActivated)
	}

	// sign successfully
	return s.jwtService.generateToken(u)
}

type SignInInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

func (i *SignInInput) Valid() error {
	if !validatePassword(i.Password) {
		return fmt.Errorf("password invalid")
	}

	return validator.New().Struct(i)
}

func (s *service) Register(input *RegisterInput) error {
	if err := validate(input); err != nil {
		return errorutil.Wrap(ErrInvalidInput, err)
	}

	if _, err := s.userRepository.FindUserByEmail(input.Email); err == nil {
		return errorutil.Wrap(ErrEmailExisted, input.Email)
	}

	if _, err := s.userRepository.FindUserByUsername(input.Username); err == nil {
		return errorutil.Wrap(ErrUsernameExisted, input.Username)
	}

	u, err := NewUser(input.Email, input.Password)
	if err != nil {
		return errorutil.Wrap(ErrUnknown, err)
	}

	id, err := s.userRepository.CreateUser(u)

	if err != nil {
		return errorutil.Wrap(ErrUnknown, err)
	}

	time.AfterFunc(time.Millisecond, func() {
		s.sender.SendActivationEmail(id)
	})

	return nil
}

type RegisterInput struct {
	Email                string `json:"email" validate:"required,email"`
	Username             string `json:"username" validate:"required,min=2,max=30"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
	FullName             string `json:"full_name" validate:"required"`
}

func (i *RegisterInput) Valid() error {
	v := validator.New()
	if !validatePassword(i.Password) {
		return fmt.Errorf("password invalid")
	}

	return v.Struct(i)
}

func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLetter := false
	hasDigit := false
	for _, r := range password {
		if unicode.IsDigit(r) {
			hasDigit = true
		}
		if unicode.IsLetter(r) {
			hasLetter = true
		}
	}

	return hasLetter && hasDigit
}
