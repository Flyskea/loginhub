package convert

import (
	"database/sql"
	"time"

	"loginhub/internal/domain/user/entity"
	"loginhub/internal/infra/persistence/user/po"
)

func UserPOToEntity(user *po.User) *entity.User {
	eu := &entity.User{
		UserID:   user.UserID,
		Name:     user.Name,
		Avatar:   user.Avatar,
		Password: user.Password,
		ActiveInfo: &entity.ActiveInfo{
			IP:          user.IP,
			LastLoginAt: time.Unix(user.LastLoginAt, 0),
		},
	}
	if user.Mobile.Valid {
		eu.Mobile = user.Mobile.String
	}
	if user.Email.Valid {
		eu.Email = user.Email.String
	}
	if user.Account.Valid {
		eu.Account = user.Account.String
	}
	return eu
}

func UserPOFromEntity(user *entity.User) *po.User {
	if user == nil {
		return nil
	}
	po := &po.User{
		UserID:      user.UserID,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Password:    user.Password,
		IP:          user.ActiveInfo.IP,
		LastLoginAt: user.ActiveInfo.LastLoginAt.Unix(),
	}

	if user.Mobile != "" {
		po.Mobile = sql.NullString{String: user.Mobile, Valid: true}
	}
	if user.Email != "" {
		po.Email = sql.NullString{String: user.Email, Valid: true}
	}
	if user.Account != "" {
		po.Account = sql.NullString{String: user.Account, Valid: true}
	}
	return po
}
