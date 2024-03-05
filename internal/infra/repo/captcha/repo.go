package captcha

import (
	"context"
	"fmt"

	"github.com/Flyskea/gotools/errors"
	"github.com/redis/go-redis/v9"

	"loginhub/internal/base/reason"
	"loginhub/internal/conf"
	"loginhub/internal/service/captcha"
)

var _ captcha.CaptchaRepository = (*CaptchaRepo)(nil)

type Captcha struct {
	Code        string
	LastRequest int64
}

const (
	CaptchaKeyPrefix = "captcha"
	userIDPrefix     = "uid"
	emailPrefix      = "email"
)

func userIDKey(t captcha.CaptchaType) string {
	return fmt.Sprintf("%s:%d:%s", CaptchaKeyPrefix, t, userIDPrefix)
}

func emailKey(t captcha.CaptchaType) string {
	return fmt.Sprintf("%s:%d:%s", CaptchaKeyPrefix, t, emailPrefix)
}

type CaptchaRepo struct {
	rdb redis.UniversalClient
}

func NewCaptchaRepo(rdb redis.UniversalClient) *CaptchaRepo {
	return &CaptchaRepo{
		rdb: rdb,
	}
}

func (r *CaptchaRepo) Save(
	ctx context.Context,
	userID int64,
	capchaType captcha.CaptchaType,
	c *captcha.Captcha,
) error {
	key := userIDKey(capchaType)
	err := r.rdb.Set(ctx, key, c.ToJSONString(), conf.GConf.Captcha.Ttl.AsDuration()).Err()
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (r *CaptchaRepo) SaveByEmail(
	ctx context.Context,
	email string,
	capchaType captcha.CaptchaType,
	c *captcha.Captcha,
) error {
	key := emailKey(capchaType)
	err := r.rdb.Set(ctx, key, c.ToJSONString(), conf.GConf.Captcha.Ttl.AsDuration()).Err()
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (r *CaptchaRepo) Get(
	ctx context.Context,
	userID int64,
	capchaType captcha.CaptchaType,
) (*captcha.Captcha, error) {
	key := userIDKey(capchaType)
	code, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.NotFound(reason.CaptchaNotExist)
		}
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return (&captcha.Captcha{}).FromJSONString(code), nil
}

func (r *CaptchaRepo) GetByEmail(
	ctx context.Context,
	email string,
	capchaType captcha.CaptchaType,
) (*captcha.Captcha, error) {
	key := emailKey(capchaType)
	code, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.NotFound(reason.CaptchaNotExist)
		}
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return (&captcha.Captcha{}).FromJSONString(code), nil
}
