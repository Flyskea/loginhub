package convert

import (
	"time"

	"loginhub/internal/domain/user/entity"
	"loginhub/internal/infra/persistence/user/po"
)

func UserPOToEntity(user *po.User) *entity.User {
	return &entity.User{
		UserID:   user.UserID,
		Name:     user.Name,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Password: user.Password,
		ActiveInfo: &entity.ActiveInfo{
			IP:          user.IP,
			LastLoginAt: time.Unix(user.LastLoginAt, 0),
		},
	}
}

func UserPOFromEntity(user *entity.User) *po.User {
	return &po.User{
		UserID:      user.UserID,
		Name:        user.Name,
		Email:       user.Email,
		Avatar:      user.Avatar,
		Password:    user.Password,
		IP:          user.ActiveInfo.IP,
		LastLoginAt: user.ActiveInfo.LastLoginAt.Unix(),
	}
}
