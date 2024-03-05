package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/Flyskea/gotools/errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/redis/go-redis/v9"

	"loginhub/internal/base/pager"
	"loginhub/internal/base/reason"
	"loginhub/internal/domain/oauth2/entity"
	"loginhub/internal/domain/oauth2/repository"
	"loginhub/internal/infra/data"
	"loginhub/internal/infra/persistence/oauth2/po"
)

var _ repository.OAuth2ProviderInfoRepository = (*OAuth2ProviderRepo)(nil)

type OAuth2ProviderRepo struct {
	txm *data.TXManager
	rdb redis.UniversalClient
}

func NewOAuth2ProviderRepo(
	txm *data.TXManager,
	rdb redis.UniversalClient,
) *OAuth2ProviderRepo {
	return &OAuth2ProviderRepo{
		txm: txm,
		rdb: rdb,
	}
}

const (
	statePrefix = "oauth2_state"
)

var (
	selectFields = []string{"id", "type", "client_id", "client_secret", "redirect_url"}
)

func stateKey(providerType, key string) string {
	return fmt.Sprintf("%s:%s:%s", statePrefix, providerType, key)
}

func (o *OAuth2ProviderRepo) GetByType(ctx context.Context, providerType string) (*entity.ProviderInfo, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sql, args := sb.Select(selectFields...).
		From(po.OAuth2ProviderTableName).
		Where(sb.EQ("type", providerType)).
		Limit(1).Build()
	po := &po.OAuth2Provider{}
	err := o.txm.DB(ctx).Raw(sql, args...).First(po).Error
	if err = data.GromToError(err); err != nil {
		return nil, err
	}
	return po.ToEntity(), nil
}

func (o *OAuth2ProviderRepo) List(
	ctx context.Context,
	pageCond *pager.PageCond,
) ([]*entity.ProviderInfo, int64, error) {
	pageCond.Validate()
	sb := sqlbuilder.NewSelectBuilder()
	sql, args := sb.Select(selectFields...).
		From(po.OAuth2ProviderTableName).
		OrderBy("id").Asc().
		Limit(int(pageCond.Size)).Offset(int(pageCond.GetOffset())).
		Build()
	pos := po.OAuth2ProviderList{}
	err := o.txm.DB(ctx).Raw(sql, args...).Scan(&pos).Error
	if err = data.GromToError(err); err != nil {
		return nil, 0, err
	}
	if len(pos) == 0 {
		return nil, 0, nil
	}

	var total int64
	sb = sqlbuilder.NewSelectBuilder()
	sql, args = sb.Select("COUNT(*)").
		From(po.OAuth2ProviderTableName).Build()
	err = o.txm.DB(ctx).Raw(sql, args...).Scan(&total).Error
	if err = data.GromToError(err); err != nil {
		return nil, 0, err
	}
	return pos.ToEntity(), total, nil
}

func (o *OAuth2ProviderRepo) Save(ctx context.Context, provider *entity.ProviderInfo) error {
	po := (&po.OAuth2Provider{}).FromEntity(provider)
	err := o.txm.DB(ctx).Save(po).Error
	if err = data.GromToError(err); err != nil {
		return err
	}
	provider.ID = po.ID
	return nil
}

func (o *OAuth2ProviderRepo) DeleteByID(ctx context.Context, id int64) error {
	sb := sqlbuilder.NewUpdateBuilder()
	sql, args := sb.Where(
		sb.Equal("id", id),
	).Set(sb.Assign("deleted_at", time.Now().Unix())).
		Update(po.OAuth2ProviderTableName).Build()

	return data.GromToError(o.txm.DB(ctx).Exec(sql, args...).Error)
}

func (o *OAuth2ProviderRepo) Update(ctx context.Context, provider *entity.ProviderInfo) error {
	sb := sqlbuilder.NewUpdateBuilder()
	if provider.ClientID != "" {
		sb.Set(sb.Assign("client_id", provider.ClientID))
	}
	if provider.ClientSecret != "" {
		sb.Set(sb.Assign("client_secret", provider.ClientSecret))
	}
	if provider.RedirectURL != "" {
		sb.Set(sb.Assign("redirect_url", provider.RedirectURL))
	}
	sql, args := sb.Where(
		sb.Equal("id", provider.ID),
	).Update(po.OAuth2ProviderTableName).Build()
	return data.GromToError(o.txm.DB(ctx).Exec(sql, args...).Error)
}

func (o *OAuth2ProviderRepo) SaveState(
	ctx context.Context,
	pType string,
	key string,
	state string,
) error {
	rk := stateKey(pType, key)
	err := o.rdb.Set(ctx, rk, state, 0).Err()
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}
func (o *OAuth2ProviderRepo) GetStateByKey(
	ctx context.Context,
	pType string,
	key string,
) (string, error) {
	rk := stateKey(pType, key)

	state, err := o.rdb.Get(ctx, rk).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", errors.NotFound(reason.ObjectNotFoundError)
		}
		return "", errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return state, nil
}
