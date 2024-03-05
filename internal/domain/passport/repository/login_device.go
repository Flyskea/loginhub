package repository

import (
	"context"

	"loginhub/internal/domain/passport/entity"
)

type LoginDeviceRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*entity.LoginDevice, error)
	GetDeviceByID(ctx context.Context, userID int64, deviceID string) (*entity.Device, error)
	Create(ctx context.Context, userID int64, device *entity.Device) error
	DeleteOneByDeviceID(ctx context.Context, userID int64, deviceID string) error
	DeleteAllByUserID(ctx context.Context, userID int64) error
}
