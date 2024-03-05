package repo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Flyskea/gotools/errors"
	"github.com/redis/go-redis/v9"

	"loginhub/internal/base/reason"
	"loginhub/internal/domain/passport/entity"
	"loginhub/internal/domain/passport/repository"
)

var _ repository.RefreshTokenRepository = (*RefreshTokenRepo)(nil)

const (
	refreshTokenPrefix = "refresh_token"

	iterCount = 30
)

func refreshTokenKey(userID int64, deviceID string, tokenID string) string {
	return fmt.Sprintf("%s:%d:%s:%s", refreshTokenPrefix, userID, deviceID, tokenID)
}

func refreshTokenKeyByEntity(t *entity.RefreshToken) string {
	return refreshTokenKey(t.UserID, t.DeviceID, t.ID)
}

type RefreshTokenRepo struct {
	rdb    redis.UniversalClient
	logger *slog.Logger
}

func NewRefreshTokenRepo(
	rdb redis.UniversalClient,
	logger *slog.Logger,
) *RefreshTokenRepo {
	return &RefreshTokenRepo{
		rdb:    rdb,
		logger: logger,
	}
}

func (r *RefreshTokenRepo) StoreToken(ctx context.Context, token *entity.RefreshToken) error {
	key := refreshTokenKeyByEntity(token)
	return r.rdb.Set(ctx, key, "", token.GetTTL()).Err()
}

func (r *RefreshTokenRepo) RotateToken(ctx context.Context, oldToken *entity.RefreshToken, newToken *entity.RefreshToken) error {
	oldKey := refreshTokenKeyByEntity(oldToken)
	newKey := refreshTokenKeyByEntity(newToken)
	cnt, err := r.rdb.Del(ctx, oldKey).Result()
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if cnt != 1 {
		r.logger.WarnContext(ctx, "refresh token not exist", slog.String("key", oldKey))
	}

	return r.rdb.Set(ctx, newKey, "", newToken.GetTTL()).Err()
}

func (r *RefreshTokenRepo) IsTokenExist(ctx context.Context, token *entity.RefreshToken) (bool, error) {
	cnt, err := r.rdb.Exists(ctx, refreshTokenKeyByEntity(token)).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return false, nil
	case err != nil:
		return false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	default:
		return cnt == 1, nil
	}
}

func (r *RefreshTokenRepo) scanAnddel(ctx context.Context, match string) error {
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = r.rdb.Scan(ctx, cursor, match, iterCount).Result()
		if err != nil {
			return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		if len(keys) == 0 {
			break
		}
		err = r.rdb.Del(ctx, keys...).Err()
		if err != nil {
			return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}

func (r *RefreshTokenRepo) RevokeTokenByUserID(ctx context.Context, userID int64) error {
	match := fmt.Sprintf("%s:%d:*", refreshTokenPrefix, userID)
	return r.scanAnddel(ctx, match)
}

func (r *RefreshTokenRepo) RevokeTokenByDeviceID(ctx context.Context, userID int64, deviceID string) error {
	match := fmt.Sprintf("%s:%d:%s:*", refreshTokenPrefix, userID, deviceID)
	return r.scanAnddel(ctx, match)
}
