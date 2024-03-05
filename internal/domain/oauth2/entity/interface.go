package entity

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type UserInfo struct {
	ID          string
	Username    string
	DisplayName string
	UnionID     string
	Email       string
	Phone       string
	CountryCode string
	AvatarUrl   string
	Extra       map[string]string
}

type Oauth2Provider interface {
	BuildRequestURL() (string, string, error)
	GetToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error)
	SetClient(client *http.Client)
}

const (
	GithubOauth2ProviderName = "github"
)
