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

type RefreshToken struct {
	Token
}

func CreateRefreshToken(
	userID int64,
	userName string,
	DeviceID string,
) *RefreshToken {
	id := convert.TrimUUID(uuid.New().String())
	return &RefreshToken{
		Token: Token{
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        id,
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(conf.GConf.GetPassport().GetRefreshTokenTtl().AsDuration())),
			},
			UserInfo: UserInfo{
				UserID:   userID,
				UserName: userName,
				DeviceID: DeviceID,
			},
		},
	}
}

func ParseAndVerifyRefreshToken(token string) (*RefreshToken, error) {
	var t RefreshToken
	jwtToken, err := jwt.ParseWithClaims(token, &t, SecretKey)
	if err != nil {
		return nil, err
	}
	if !jwtToken.Valid {
		return nil, err
	}
	return &t, nil
}

func (t *RefreshToken) Sign() (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, t).SignedString([]byte(conf.GConf.GetPassport().GetSecret()))
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	return token, nil
}

func (t *RefreshToken) GetTTL() time.Duration {
	return conf.GConf.GetPassport().GetRefreshTokenTtl().AsDuration()
}

func (t *RefreshToken) Refresh() *RefreshToken {
	return CreateRefreshToken(t.UserID, t.UserName, t.DeviceID)
}
