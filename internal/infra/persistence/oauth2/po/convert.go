package po

import (
	"golang.org/x/oauth2"

	"loginhub/internal/domain/oauth2/entity"
)

func (p *OAuth2Provider) ToEntity() *entity.ProviderInfo {
	return &entity.ProviderInfo{
		ID:           p.ID,
		Type:         p.Type,
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		RedirectURL:  p.RedirectUrl,
	}
}

func (p *OAuth2Provider) FromEntity(provider *entity.ProviderInfo) *OAuth2Provider {
	p.ID = provider.ID
	p.Type = provider.Type
	p.ClientID = provider.ClientID
	p.ClientSecret = provider.ClientSecret
	p.RedirectUrl = provider.RedirectURL

	return p
}

func (l OAuth2ProviderList) ToEntity() []*entity.ProviderInfo {
	providers := make([]*entity.ProviderInfo, 0, len(l))
	for _, p := range l {
		providers = append(providers, p.ToEntity())
	}
	return providers
}

func (u *UserThirdAuth) FromEntity(
	userID int64,
	providerType string,
	token *oauth2.Token,
	info *entity.UserInfo,
) *UserThirdAuth {
	u.UserID = userID
	u.Type = providerType
	u.AuthID = info.ID
	u.UnionID = info.UnionID
	u.Credential = token.AccessToken
	u.RefreshToken = token.RefreshToken
	return u
}
