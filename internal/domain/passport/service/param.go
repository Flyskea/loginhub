package service

import (
	passportentity "loginhub/internal/domain/passport/entity"
	"loginhub/internal/domain/user/entity"
)

type LoginInfo struct {
	Email    string
	Phone    string
	Password string
	Device   *passportentity.Device
}

type AuthenticationInfo struct {
	AccessToken  *passportentity.AccessToken
	RefreshToken *passportentity.RefreshToken
	User         *entity.User
}

type RegisterType uint8

const (
	EmailRegisterType = iota + 1
)

type RegisterInfo struct {
	RegisterType RegisterType
	User         *entity.User
	Device       *passportentity.Device
}
