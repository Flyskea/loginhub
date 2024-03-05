package entity

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github"

	"loginhub/pkg/convert"
)

var _ Oauth2Provider = (*GithubOauth2Provider)(nil)

type GithubOauth2Provider struct {
	BasicOauth2Provider
}

type GitHubUserInfo struct {
	Login                   string    `json:"login"`
	ID                      int       `json:"id"`
	NodeID                  string    `json:"node_id"`
	AvatarURL               string    `json:"avatar_url"`
	GravatarID              string    `json:"gravatar_id"`
	URL                     string    `json:"url"`
	HTMLURL                 string    `json:"html_url"`
	FollowersURL            string    `json:"followers_url"`
	FollowingURL            string    `json:"following_url"`
	GistsURL                string    `json:"gists_url"`
	StarredURL              string    `json:"starred_url"`
	SubscriptionsURL        string    `json:"subscriptions_url"`
	OrganizationsURL        string    `json:"organizations_url"`
	ReposURL                string    `json:"repos_url"`
	EventsURL               string    `json:"events_url"`
	ReceivedEventsURL       string    `json:"received_events_url"`
	Type                    string    `json:"type"`
	SiteAdmin               bool      `json:"site_admin"`
	Name                    string    `json:"name"`
	Company                 string    `json:"company"`
	Blog                    string    `json:"blog"`
	Location                string    `json:"location"`
	Email                   string    `json:"email"`
	Hireable                bool      `json:"hireable"`
	Bio                     string    `json:"bio"`
	TwitterUsername         string    `json:"twitter_username"`
	PublicRepos             int       `json:"public_repos"`
	PublicGists             int       `json:"public_gists"`
	Followers               int       `json:"followers"`
	Following               int       `json:"following"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	PrivateGists            int       `json:"private_gists"`
	TotalPrivateRepos       int       `json:"total_private_repos"`
	OwnedPrivateRepos       int       `json:"owned_private_repos"`
	DiskUsage               int       `json:"disk_usage"`
	Collaborators           int       `json:"collaborators"`
	TwoFactorAuthentication bool      `json:"two_factor_authentication"`
	Plan                    struct {
		Name          string `json:"name"`
		Space         int    `json:"space"`
		PrivateRepos  int    `json:"private_repos"`
		Collaborators int    `json:"collaborators"`
	} `json:"plan"`
}

func NewGithubOauth2Provider(
	clientID string,
	clientSecret string,
	redirectURL string,
) *GithubOauth2Provider {
	p := GithubOauth2Provider{
		BasicOauth2Provider: BasicOauth2Provider{
			Provider: GithubOauth2ProviderName,
		},
	}

	config := p.getConfig()
	config.ClientID = clientID
	config.ClientSecret = clientSecret
	config.RedirectURL = redirectURL

	p.Config = config
	return &p
}

func (p *GithubOauth2Provider) getConfig() *oauth2.Config {
	return &oauth2.Config{
		Scopes:   []string{"user:email", "read:user"},
		Endpoint: githubOAuth.Endpoint,
	}
}

func (p *GithubOauth2Provider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var githubUserInfo GitHubUserInfo
	err = json.Unmarshal(body, &githubUserInfo)
	if err != nil {
		return nil, err
	}
	slog.InfoContext(ctx, "github user info", slog.String("user", convert.BytesToString(body)))
	userInfo := UserInfo{
		ID:          strconv.Itoa(githubUserInfo.ID),
		Username:    githubUserInfo.Login,
		DisplayName: githubUserInfo.Name,
		Email:       githubUserInfo.Email,
		AvatarUrl:   githubUserInfo.AvatarURL,
	}
	return &userInfo, nil
}
