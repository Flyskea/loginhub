package repository

import (
	"context"
	"loginhub/internal/domain/user/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByUserID(ctx context.Context, userID int64) (*entity.User, error)
	GetByMobile(ctx context.Context, mobile string) (*entity.User, error)
	GetByAccountName(ctx context.Context, accountName string) (*entity.User, error)
}
