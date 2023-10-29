package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/iwanhae/oidc-go-demo/pkg/errors"

	"golang.org/x/oauth2"
)

type Config struct {
	Provider     string `env:"OIDC_PROVIDER" envDefault:"https://accounts.google.com"`
	ClientID     string `env:"OIDC_CLIENT_ID,required"`
	ClientSecret string `env:"OIDC_CLIENT_SECRET,required"`
	RedirectURL  string `env:"OIDC_REDIRECT_URL,required"`
}

type OIDCService interface {
	GetRedirectURL() string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	FetchUserInfo(ctx context.Context, token *oauth2.Token) (*oidc.UserInfo, error)
}

type oidcService struct {
	*oauth2.Config
	provider *oidc.Provider
}

func NewOIDCService(ctx context.Context, cfg *Config) (OIDCService, error) {
	provider, err := oidc.NewProvider(ctx, cfg.Provider)
	if err != nil {
		return nil, errors.Wrap(err, "oidc: failed to create provider")
	}
	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),

		// https://auth0.com/docs/get-started/apis/scopes/openid-connect-scopes
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}
	return &oidcService{
		Config:   &oauth2Config,
		provider: provider,
	}, nil
}

func (s *oidcService) GetRedirectURL() string {
	state := make([]byte, 16)
	_, err := rand.Read(state)
	if err != nil {
		panic(errors.Wrap(err, "oidc: failed while generating random state"))
	}
	return s.AuthCodeURL(base64.RawURLEncoding.EncodeToString(state), oauth2.AccessTypeOffline)
}

func (s *oidcService) FetchUserInfo(ctx context.Context, token *oauth2.Token) (*oidc.UserInfo, error) {
	return s.provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
}
