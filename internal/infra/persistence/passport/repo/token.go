package repo

import (
	"context"
	"fmt"

	"github.com/Flyskea/gotools/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"loginhub/internal/base/reason"
	"loginhub/internal/domain/passport/entity"
	"loginhub/internal/domain/passport/repository"
)

const (
	accessTokenPrefix = "access_token"
)

var _ repository.AccessTokenRepository = (*AccessTokenRepo)(nil)

func accessTokenKey(userID int64, deviceID string, tokenID string) string {
	return fmt.Sprintf("%s:%d:%s:%s", accessTokenPrefix, userID, deviceID, tokenID)
}

func accessTokenKeyByEntity(t *entity.AccessToken) string {
	return accessTokenKey(t.UserID, t.DeviceID, t.ID)
}

type AccessTokenRepo struct {
	db  *gorm.DB
	rdb redis.UniversalClient
}

func NewAccessTokenRepo(
	db *gorm.DB,
	rdb redis.UniversalClient,
) *AccessTokenRepo {
	return &AccessTokenRepo{
		db:  db,
		rdb: rdb,
	}
}

func (r *AccessTokenRepo) StoreToken(ctx context.Context, token *entity.AccessToken) error {
	key := accessTokenKeyByEntity(token)
	return r.rdb.Set(ctx, key, "", token.GetTTL()).Err()
}

func (r *AccessTokenRepo) IsTokenExist(ctx context.Context, token *entity.AccessToken) (bool, error) {
	cnt, err := r.rdb.Exists(ctx, accessTokenKeyByEntity(token)).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return false, nil
	case err != nil:
		return false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	default:
		return cnt == 1, nil
	}
}

func (r *AccessTokenRepo) scanAnddel(ctx context.Context, match string) error {
	fmt.Println(match)
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

func (r *AccessTokenRepo) RevokeTokenByUserID(ctx context.Context, userID int64) error {
	match := fmt.Sprintf("%s:%d:*", accessTokenPrefix, userID)
	return r.scanAnddel(ctx, match)
}

func (r *AccessTokenRepo) RevokeTokenByDeviceID(ctx context.Context, userID int64, deviceID string) error {
	match := fmt.Sprintf("%s:%d:%s:*", accessTokenPrefix, userID, deviceID)
	return r.scanAnddel(ctx, match)
}
