package auth

import (
	"context"
	"errors"
	"github.com/victornm/es-backend/pkg/auth/internal"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOAuth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"strings"
)

const (
	ProviderGoogle string = "google"
)

var ErrInvalidOAuth2Provider = errors.New("oauth2 provider not supported")

type OAuth2RegisterService interface {
	OAuth2Register(provider string) (string, error)
	OAuth2RegisterCallback(state, code string) error
}

type SetClient interface {
	setClient(factory OAuth2ClientFactory)
}

type setConfig struct {
	provider string
	client   OAuth2Client
}

func (setter *setConfig) setClient(factory OAuth2ClientFactory) {
	factory[strings.ToLower(setter.provider)] = setter.client
}

func WithGoogle(clientID, clientSecret, redirectURL string) *setConfig {
	return &setConfig{
		provider: ProviderGoogle,
		client: &googleClient{
			config: oauth2.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Endpoint:     google.Endpoint,
				RedirectURL:  redirectURL,
				Scopes: []string{
					googleOAuth2.UserinfoEmailScope,
					googleOAuth2.UserinfoProfileScope,
				},
			},
		},
	}
}

func NewOAuth2ClientFactory(options ...SetClient) OAuth2ClientFactory {
	f := OAuth2ClientFactory{}

	for _, o := range options {
		o.setClient(f)
	}

	return f
}

type OAuth2ClientFactory map[string]OAuth2Client

func (f OAuth2ClientFactory) GetService(provider string) (OAuth2Client, error) {
	if s, ok := f[strings.ToLower(provider)]; ok {
		return s, nil
	}

	return nil, ErrInvalidOAuth2Provider
}

type OAuth2Client interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	GetUser(code string) (*internal.User, error)
}

type oauth2RegisterService struct {
	stateRepository OAuth2StateRepository
	userRepository  UserRepository
	factory         OAuth2ClientFactory
}

func (s *oauth2RegisterService) OAuth2Register(provider string) (string, error) {
	oauth2Service, err := s.factory.GetService(provider)
	if err != nil {
		return "", err
	}

	state := internal.NewOAuth2State(provider)
	err = s.stateRepository.CreateState(state)
	if err != nil {
		return "", err
	}

	return oauth2Service.AuthCodeURL(state.Nonce, oauth2.AccessTypeOffline), nil
}

type googleClient struct {
	config oauth2.Config
}

func (r *googleClient) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return r.config.AuthCodeURL(state, opts...)
}

func (r *googleClient) GetUser(code string) (*internal.User, error) {
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

	u := internal.NewOAuth2User(googleUser.Email, ProviderGoogle)
	if len(googleUser.Name) > 0 {
		u.FullName = googleUser.Name
	} else {
		u.FullName = googleUser.GivenName + " " + googleUser.FamilyName
	}

	return u, nil
}

func (s *oauth2RegisterService) OAuth2RegisterCallback(nonce, code string) error {
	state, err := s.stateRepository.FindState(nonce)
	if err != nil {
		return wrapError(ErrNotAuthenticated, "invalid nonce")
	}

	oauth2Service, err := s.factory.GetService(state.Provider)
	if err != nil {
		return err
	}

	u, err := oauth2Service.GetUser(code)
	if err != nil {
		return wrapError(ErrNotAuthenticated)
	}

	if _, err := s.userRepository.FindUserByEmail(u.Email); err != nil {
		return wrapError(ErrEmailExisted, u.Email)
	}

	if _, err := s.userRepository.CreateUser(u); err != nil {
		return wrapError(ErrUnknown, "create user failed")
	}

	return nil
}

func NewOAuth2RegisterService(stateRepository OAuth2StateRepository, userRepository UserRepository, factory OAuth2ClientFactory) *oauth2RegisterService {
	return &oauth2RegisterService{
		stateRepository: stateRepository,
		factory:         factory,
		userRepository:  userRepository,
	}
}
