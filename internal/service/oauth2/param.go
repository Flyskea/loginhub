package oauth2

import (
	"golang.org/x/oauth2"

	"loginhub/internal/domain/oauth2/entity"
)

type GetUserInfoAction struct {
	ProviderType    string
	Code            string
	RequestState    string
	OAuth2SessionID string
}

type GetUserInfoResult struct {
	Token    *oauth2.Token
	UserInfo *entity.UserInfo
}
