package entity

import (
	"context"
	"net/http"

	"github.com/Flyskea/gotools/errors"
	"golang.org/x/oauth2"

	"loginhub/internal/base/reason"
	"loginhub/pkg/random"
)

var _ Oauth2Provider = (*BasicOauth2Provider)(nil)

type BasicOauth2Provider struct {
	ID       int64
	Provider string
	Config   *oauth2.Config
	Client   *http.Client
}

const (
	stateLength = 16
)

func (o *BasicOauth2Provider) generateState() string {
	return random.RandomString(random.Alphanumeric, stateLength)
}

func (o *BasicOauth2Provider) BuildRequestURL() (string, string, error) {
	state := o.generateState()

	return o.Config.AuthCodeURL(state), state, nil
}

func (o *BasicOauth2Provider) GetToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return o.Config.Exchange(ctx, code)
}

func (o *BasicOauth2Provider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	return nil, nil
}

func (o *BasicOauth2Provider) SetClient(client *http.Client) {
	o.Client = client
}

func NewOauth2Provider(
	provider string,
	clientID string,
	clientSecret string,
	redirectURL string,
) (Oauth2Provider, error) {
	if redirectURL == "" {
		return nil, errors.BadRequest(reason.OAuth2RedirectURLEmpty)
	}
	if clientID == "" {
		return nil, errors.BadRequest(reason.OAuth2ClientIDEmpty)
	}
	if clientSecret == "" {
		return nil, errors.BadRequest(reason.OAuth2ClientSecretEmpty)
	}

	switch provider {
	case GithubOauth2ProviderName:
		return NewGithubOauth2Provider(clientID, clientSecret, redirectURL), nil
	default:
		return nil, errors.BadRequest(reason.OAuth2ProviderNotImplemented)
	}
}

var SupportedProviders = []string{
	GithubOauth2ProviderName,
}

func GetSupportedProviders() []string {
	return SupportedProviders
}
