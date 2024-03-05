package repository

import (
	"context"

	"golang.org/x/oauth2"

	"loginhub/internal/base/pager"
	"loginhub/internal/domain/oauth2/entity"
	userentity "loginhub/internal/domain/user/entity"
)

type OAuth2ProviderInfoRepository interface {
	DeleteByID(ctx context.Context, id int64) error
	GetByType(ctx context.Context, providerType string) (*entity.ProviderInfo, error)
	List(ctx context.Context, pageCond *pager.PageCond) ([]*entity.ProviderInfo, int64, error)
	Save(ctx context.Context, provider *entity.ProviderInfo) error
	Update(ctx context.Context, provider *entity.ProviderInfo) error
	SaveState(ctx context.Context, providerType, key, state string) error
	GetStateByKey(ctx context.Context, providerType, key string) (string, error)
}

type OAuth2UserRepository interface {
	Save(ctx context.Context, userID int64, providerType string, token *oauth2.Token, info *entity.UserInfo) error
	GetByOAuth2UserID(ctx context.Context, oauth2UserID string) (*userentity.User, error)
}
