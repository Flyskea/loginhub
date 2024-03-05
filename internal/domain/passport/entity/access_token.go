package entity

import (
	"time"

	"github.com/Flyskea/gotools/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"loginhub/internal/base/reason"
	"loginhub/internal/conf"
	"loginhub/pkg/convert"
)

type AccessToken struct {
	Token
}

func CreateAccessToken(
	userID int64,
	userName string,
	DeviceID string,
) *AccessToken {
	id := convert.TrimUUID(uuid.New().String())
	return &AccessToken{
		Token: Token{
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        id,
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(conf.GConf.GetPassport().GetAccessTokenTtl().AsDuration())),
			},
			UserInfo: UserInfo{
				UserID:   userID,
				UserName: userName,
				DeviceID: DeviceID,
			},
		},
	}
}

func ParseAndVerifyAccessToken(token string) (*AccessToken, error) {
	var t AccessToken
	jwtToken, err := jwt.ParseWithClaims(token, &t, SecretKey)
	if err != nil {
		return nil, err
	}
	if !jwtToken.Valid {
		return nil, err
	}
	return &t, nil
}

func (t *AccessToken) Sign() (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, t).SignedString([]byte(conf.GConf.GetPassport().GetSecret()))
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	return token, nil
}

func (t *AccessToken) GetTTL() time.Duration {
	return conf.GConf.GetPassport().GetAccessTokenTtl().AsDuration()
}
