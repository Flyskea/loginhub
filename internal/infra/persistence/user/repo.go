package user

import (
	"context"
	"log/slog"

	"github.com/Flyskea/gotools/errors"
	"github.com/huandu/go-sqlbuilder"
	"gorm.io/gorm"

	"loginhub/internal/base/reason"
	"loginhub/internal/domain/user/entity"
	"loginhub/internal/domain/user/repository"
	"loginhub/internal/infra/data"
	"loginhub/internal/infra/persistence/user/convert"
	"loginhub/internal/infra/persistence/user/po"
)

type UserRepo struct {
	txm    *data.TXManager
	logger *slog.Logger
}

func New(
	txm *data.TXManager,
	logger *slog.Logger,
) *UserRepo {
	return &UserRepo{
		txm:    txm,
		logger: logger,
	}
}

var _ repository.UserRepository = (*UserRepo)(nil)

var allColumns = []string{"user_id", "name", "password", "email",
	"mobile", "account", "avatar", "ip", "last_login_at"}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	err := data.GromToError(r.txm.DB(ctx).Create(convert.UserPOFromEntity(user)).Error)
	if err != nil {
		if errors.IsConflict(err) {
			return errors.BadRequest(reason.UserNotExist)
		}
		return err
	}
	return nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sql, args := sb.Select(allColumns...).
		From("user").
		Where(sb.EQ("email", email),
			sb.EQ("deleted_at", 0)).
		Limit(1).Build()
	userPO := po.User{}
	err := r.txm.DB(ctx).Raw(sql, args...).
		First(&userPO).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound(reason.UserNotExist)
		}
		return nil, err
	}
	return convert.UserPOToEntity(&userPO), nil
}

func (r *UserRepo) GetByUserID(ctx context.Context, userID int64) (*entity.User, error) {
	sb := sqlbuilder.NewSelectBuilder()
	userPO := po.User{}
	sql, arg := sb.Select(allColumns...).
		From("user").
		Where(sb.EQ("user_id", userID),
			sb.EQ("deleted_at", 0)).
		Limit(1).Build()
	err := r.txm.DB(ctx).Raw(sql, arg...).First(&userPO).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound(reason.UserNotExist)
		}
		return nil, err
	}
	return convert.UserPOToEntity(&userPO), nil
}

func (r *UserRepo) GetByMobile(ctx context.Context, mobile string) (*entity.User, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sql, arg := sb.Select(allColumns...).
		From("user").
		Where(sb.EQ("mobile", mobile),
			sb.EQ("deleted_at", 0)).
		Limit(1).Build()
	userPO := po.User{}
	err := r.txm.DB(ctx).Raw(sql, arg...).First(&userPO).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound(reason.UserNotExist)
		}
		return nil, err
	}
	return convert.UserPOToEntity(&userPO), nil
}

func (r *UserRepo) GetByAccountName(ctx context.Context, accountName string) (*entity.User, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sql, arg := sb.Select(allColumns...).
		From("user").
		Where(sb.EQ("account", accountName),
			sb.EQ("deleted_at", 0)).
		Limit(1).Build()
	userPO := po.User{}
	err := r.txm.DB(ctx).Raw(sql, arg...).First(&userPO).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound(reason.UserNotExist)
		}
		return nil, err
	}
	return convert.UserPOToEntity(&userPO), nil
}
