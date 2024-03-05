package repository

import (
	"context"

	"loginhub/internal/domain/passport/entity"
)

type AccessTokenRepository interface {
	StoreToken(ctx context.Context, token *entity.AccessToken) error
	IsTokenExist(ctx context.Context, token *entity.AccessToken) (bool, error)
	RevokeTokenByUserID(ctx context.Context, userID int64) error
	RevokeTokenByDeviceID(ctx context.Context, userID int64, deviceID string) error
}
