package passport

import (
	passportentity "loginhub/internal/domain/passport/entity"
	"loginhub/internal/domain/user/entity"
	"loginhub/internal/service/captcha"
)

type LoginResult struct {
	AccessToken    string
	AccessTokenTTL int
	RefreshToken   string
	User           *entity.User
}

type RefreshResult struct {
	AccessToken    string
	AccessTokenTTL int
	RefreshToken   string
	User           *entity.User
}

type DeleteDevice struct {
	DeviceID string
	UserInfo passportentity.UserInfo
}

type RegisterType uint8

const (
	EmailRegisterType = iota + 1
)

type RegisterInfo struct {
	RegisterType    RegisterType
	UserName        string
	Email           string
	Phone           string
	Password        string
	PasswordConfirm string
	Code            string
	Device          *passportentity.Device
}

type CaptchaType uint8

const (
	RegisterCaptchatType = iota + 1
)

func (t CaptchaType) ToCaptchaService() captcha.CaptchaType {
	switch t {
	case RegisterCaptchatType:
		return captcha.RegisterCaptchatType
	}
	return 0
}

type EmailSend struct {
	Email       string
	CaptchaType CaptchaType
}
