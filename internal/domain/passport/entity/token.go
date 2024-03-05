package entity

import (
	"loginhub/internal/conf"
	"loginhub/pkg/convert"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	jwt.RegisteredClaims
	UserInfo
}

func SecretKey(*jwt.Token) (interface{}, error) {
	return convert.StringToBytes(conf.GConf.GetPassport().GetSecret()), nil
}

func (t *Token) GetUserInfo() UserInfo {
	return t.UserInfo
}
