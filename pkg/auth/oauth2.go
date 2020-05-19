package auth

import (
	"context"
	"strings"

	"github.com/victornm/es-backend/pkg/errorutil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOAuth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type OAuth2Service interface {
	OAuth2Register(input OAuth2Input) error
	OAuth2SignIn(input OAuth2Input) (string, error)
}

type OAuth2Config struct {
	UserRepository UserRepository
	JWTService     JWTService

	Providers []OAuth2Provider
}

func NewOAuth2Service(config *OAuth2Config) OAuth2Service {
	return &oauth2Service{
		userRepository: config.UserRepository,
		factory:        newProviderFactory(config.Providers...),
		jwtService:     config.JWTService,
	}
}

type OAuth2Provider interface {
	Name() string
	GetUser(code string) (*User, error)
}

type oauth2Service struct {
	userRepository UserRepository
	factory        providerFactory
	jwtService     JWTService
}

type OAuth2Input struct {
	Provider string `json:"provider"`
	Code     string `json:"code"`
}

func (s *oauth2Service) OAuth2Register(input OAuth2Input) error {
	client, ok := s.factory.getProvider(input.Provider)
	if !ok {
		return errorutil.Wrap(ErrInvalidOAuth2Provider)
	}

	u, err := client.GetUser(input.Code)
	if err != nil {
		return errorutil.Wrap(ErrNotAuthenticated, err)
	}

	_, err = s.userRepository.FindUserByEmail(u.Email)
	if err == nil {
		return errorutil.Wrap(ErrEmailExisted, err)
	}

	_, err = s.userRepository.CreateUser(u)
	if err != nil {
		return errorutil.Wrap(ErrUnknown, err)
	}

	return nil
}

func (s *oauth2Service) OAuth2SignIn(input OAuth2Input) (string, error) {
	client, ok := s.factory.getProvider(input.Provider)
	if !ok {
		return "", errorutil.Wrap(ErrInvalidOAuth2Provider)
	}

	u, err := client.GetUser(input.Code)
	if err != nil {
		return "", errorutil.Wrap(ErrNotAuthenticated, err)
	}

	user, err := s.userRepository.FindUserByEmail(u.Email)
	if err != nil {
		return "", errorutil.Wrap(ErrNotAuthenticated, err)
	}

	if user.Provider != u.Provider {
		return "", errorutil.Wrap(ErrNotAuthenticated, "provider is not the same")
	}

	if !user.IsActive {
		return "", errorutil.Wrap(ErrNotActivated)
	}

	return s.jwtService.generateToken(user)
}

func newProviderFactory(providers ...OAuth2Provider) providerFactory {
	f := providerFactory{}

	for _, p := range providers {
		f.setProvider(p)
	}

	return f
}

type providerFactory map[string]OAuth2Provider

func (f providerFactory) getProvider(provider string) (OAuth2Provider, bool) {
	if s, ok := f[strings.ToLower(provider)]; ok {
		return s, true
	}

	return nil, false
}

func (f providerFactory) setProvider(provider OAuth2Provider) {
	f[strings.ToLower(provider.Name())] = provider
}

func NewGoogleProvider(clientID, clientSecret string) *googleProvider {
	return &googleProvider{
		config: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			Scopes: []string{
				googleOAuth2.UserinfoEmailScope,
				googleOAuth2.UserinfoProfileScope,
			},
		},
	}
}

type googleProvider struct {
	config oauth2.Config
}

func (r *googleProvider) Name() string {
	return "google"
}

func (r *googleProvider) GetUser(code string) (*User, error) {
	ctx := context.Background()

	token, err := r.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	googleService, err := googleOAuth2.NewService(ctx, option.WithTokenSource(r.config.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}

	srv := googleOAuth2.NewUserinfoV2Service(googleService)
	googleUser, err := srv.Me.Get().Do()
	if err != nil {
		return nil, err
	}

	u := NewOAuth2User(googleUser.Email, r.Name())
	if len(googleUser.Name) > 0 {
		u.FullName = googleUser.Name
	} else {
		u.FullName = googleUser.GivenName + " " + googleUser.FamilyName
	}

	return u, nil
}
