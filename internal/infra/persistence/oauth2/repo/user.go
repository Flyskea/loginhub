package repo

import (
	"context"

	"github.com/huandu/go-sqlbuilder"
	"golang.org/x/oauth2"

	"loginhub/internal/domain/oauth2/entity"
	"loginhub/internal/domain/oauth2/repository"
	userentity "loginhub/internal/domain/user/entity"
	"loginhub/internal/infra/data"
	"loginhub/internal/infra/persistence/oauth2/po"
	userconvert "loginhub/internal/infra/persistence/user/convert"
	userpo "loginhub/internal/infra/persistence/user/po"
)

var _ repository.OAuth2UserRepository = (*UserThirdAuthRepo)(nil)

var userAllColumns = []string{"user_id", "name", "password", "email",
	"mobile", "account", "avatar", "ip", "last_login_at"}

type UserThirdAuthRepo struct {
	txm *data.TXManager
}

func NewUserThirdAuthRepo(
	txm *data.TXManager,
) *UserThirdAuthRepo {
	return &UserThirdAuthRepo{
		txm: txm,
	}
}

func (r *UserThirdAuthRepo) Save(
	ctx context.Context,
	userID int64,
	providerType string,
	token *oauth2.Token,
	info *entity.UserInfo,
) error {
	uta := (&po.UserThirdAuth{}).FromEntity(userID, providerType, token, info)
	return data.GromToError(r.txm.DB(ctx).Create(uta).Error)
}

func (r *UserThirdAuthRepo) GetByOAuth2UserID(ctx context.Context, oauth2UserID string) (*userentity.User, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sql, args := sb.Select(userAllColumns...).
		From(userpo.UserTableName).
		JoinWithOption(sqlbuilder.InnerJoin,
			sb.As(po.UserThirdAuthTableName, "uta"),
			sb.EQ("uta.user_id", "user.user_id"),
			sb.EQ("uta.oauth2_user_id", oauth2UserID)).
		Limit(1).Build()

	userPO := userpo.User{}
	err := data.GromToError(r.txm.DB(ctx).Raw(sql, args...).Scan(&userPO).Error)
	if err != nil {
		return nil, err
	}
	return userconvert.UserPOToEntity(&userPO), nil
}
