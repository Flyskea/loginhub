package passport

import (
	"time"

	"loginhub/internal/service/passport"
)

type LoginType uint8

const (
	UnknownLoginType = iota
	LocalPasswordLoginType
	OAuth2LoginType
)

func (t LoginType) ToService() passport.LoginType {
	switch t {
	case LocalPasswordLoginType:
		return passport.LocalPasswordLoginType
	case OAuth2LoginType:
		return passport.OAuth2LoginType
	default:
		return passport.UnknownLoginType
	}
}

type LocalPasswordLogin struct {
	Account  string `json:"account" binding:"required_if=LoginType 1"`
	Password string `json:"password" binding:"required_if=LoginType 1"`
}

type Oauth2Login struct {
	Provider string `json:"provider" binding:"required_if=LoginType 2"`
	State    string `json:"state" binding:"required_if=LoginType 2"`
	Code     string `json:"code" binding:"required_if=LoginType 2"`
}

type LoginRequest struct {
	LoginType LoginType `json:"type" binding:"required,gte=1"`
	LocalPasswordLogin
	Oauth2Login
}

type User struct {
	UserID int64  `json:"uid"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type LoginResponse struct {
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

type RegisterType uint8

const (
	EmailRegisterType = iota + 1
)

type RegisterRequest struct {
	RegisterType    RegisterType `json:"register_type" binding:"required,gte=1"`
	Email           string       `json:"email" binding:"required_if=RegisterType 1,emailVerify" label:"邮箱"`
	Password        string       `json:"password" binding:"required"`
	PasswordConfirm string       `json:"password_confirm" binding:"required"`
	UserName        string       `json:"user_name" binding:"required"`
	Captcha         string       `json:"capcha" binding:"required"`
}

type RegisterResponse struct {
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

type EmailSendType uint8

func (t EmailSendType) ToService() passport.CaptchaType {
	switch t {
	case RegisterEmailSendType:
		return passport.RegisterCaptchatType
	}
	return 0
}

const (
	RegisterEmailSendType EmailSendType = iota + 1 // 注册验证码
)

type EmailSendRequest struct {
	Email string        `json:"email" binding:"required,emailVerify"`
	Type  EmailSendType `json:"type" binding:"required"`
}

type SessionRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type SessionRefreshResponse struct {
	RefreshToken string `json:"refresh_token"`
}

type Device struct {
	ID        string
	OS        string
	Browser   string
	IP        string
	Location  string
	CreatedAt time.Time
}

type LoginDevice struct {
	Devices []*Device `json:"devices"`
}

type KickDeviceRequest struct {
	DeviceID string `uri:"id" binding:"required"`
}
