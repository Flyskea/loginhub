package oauth2

import (
	"loginhub/internal/base/pager"
	"loginhub/internal/domain/oauth2/entity"
)

const (
	OAuth2SessionName = "oauth2_session"
)

type ProviderURI struct {
	ID int64 `uri:"id" binding:"required"`
}

type GetOauthRequestURLRequest struct {
	Provider string `uri:"provider" binding:"required"`
}

type GetOauthRequestURLResponse struct {
	RedirectURL string `json:"redirect_url"`
}

type CreateProviderRequest struct {
	Type         string `json:"type" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	RedirectURL  string `json:"redirect_url" binding:"required,url"`
}

func (c *CreateProviderRequest) ToProviderInfo() *entity.ProviderInfo {
	return &entity.ProviderInfo{
		Type:         c.Type,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
	}
}

type CreateProviderResponse struct {
	entity.ProviderInfo
}

type ListProviderRequest struct {
	pager.PageCond
}

type ListProviderResponse struct {
	pager.PageResp[*entity.ProviderInfo]
}

type UpdateProviderRequest struct {
	CreateProviderRequest
}

type UpdateProviderResponse struct {
	entity.ProviderInfo
}

type DeleteProviderRequest struct {
	ProviderURI
}

type GetSupportedProvidersResponse struct {
	Providers []string `json:"providers"`
}
