package auth

import "errors"

var (
	// Authentication errors
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrNotActivated     = errors.New("not activated")

	// Registration errors
	ErrEmailExisted    = errors.New("email already existed")
	ErrUsernameExisted = errors.New("username already existed")

	// OAuth2 errors
	ErrInvalidOAuth2Provider = errors.New("oauth2 provider not supported")

	// Common errors
	ErrInvalidInput = errors.New("invalid input")
	ErrUnknown      = errors.New("unknown error")
)
