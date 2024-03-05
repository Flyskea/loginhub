package passport

import (
	"time"

	"loginhub/internal/service/passport"
)

type LoginType uint8

const (
	UnknownLoginType = iota
	EmailLoginType
	PhoneLoginType
	AccountLoginType
	Oauth2LoginType
)

type Oauth2Login struct {
	Provider string `json:"provider" binding:"required_if=LoginType 4"`
	State    string `json:"state" binding:"required_if=LoginType 4"`
	Code     string `json:"code" binding:"required_if=LoginType 4"`
}

type LoginRequest struct {
	LoginType LoginType `json:"login_type" binding:"required,gte=1"`
	Email     string    `json:"email" binding:"required_if=LoginType 1,emailVerify"`
	Phone     string    `json:"phone" binding:"required_if=LoginType 2"`
	Password  string    `json:"password" binding:"required_if=LoginType 1,required_if=LoginType 2,required_if=LoginType 3"`
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

type GetOauthRedirectURLRequest struct {
	ProviderURI
}

type GetOauthRedirectURLResponse struct {
	RedirectURL string `json:"redirect_url"`
}

type CreateOauth2ProviderRequest struct {
	Provider     string `json:"provider" binding:"required"`
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	Redirect_url string `json:"redirect_url" binding:"required"`
}

type ProviderURI struct {
	Provider string `uri:"provider" binding:"required"`
}

type UpdateOauth2ProviderRequest struct {
	Provider     string `json:"-"`
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	Redirect_url string `json:"redirect_url" binding:"required"`
}

type DeleteOauth2ProviderRequest struct {
	ProviderURI
}
