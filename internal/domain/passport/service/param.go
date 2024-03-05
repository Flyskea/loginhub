package service

import (
	passportentity "loginhub/internal/domain/passport/entity"
	"loginhub/internal/domain/user/entity"
)

type AuthenticationParams struct {
	Password string // if empty, do not compare password, eg. oauth2 login
	User     *entity.User
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
