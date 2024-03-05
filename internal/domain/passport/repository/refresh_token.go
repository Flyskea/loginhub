package repository

import (
	"context"

	"loginhub/internal/domain/passport/entity"
)

type RefreshTokenRepository interface {
	StoreToken(ctx context.Context, token *entity.RefreshToken) error
	RotateToken(ctx context.Context, oldToken, newToken *entity.RefreshToken) error
	IsTokenExist(ctx context.Context, token *entity.RefreshToken) (bool, error)
	RevokeTokenByUserID(ctx context.Context, userID int64) error
	RevokeTokenByDeviceID(ctx context.Context, userID int64, deviceID string) error
}
