package authn

import (
	"context"
	"net/http"

	"github.com/Flyskea/gotools/errors"
	"github.com/gin-gonic/gin"

	"loginhub/internal/base/handler"
	"loginhub/internal/base/reason"
	"loginhub/internal/domain/passport/entity"
	"loginhub/internal/service/passport"
)

type ctxUserInfoKey struct{}

var ctxUserInfo = ctxUserInfoKey{}

const (
	AuthNHTTPCookieName = "mystic"
)

type AuthNUserMiddleware struct {
	ps *passport.PassportService
}

func NewAuthNUserMiddleware(ps *passport.PassportService) *AuthNUserMiddleware {
	return &AuthNUserMiddleware{
		ps: ps,
	}
}

func (m *AuthNUserMiddleware) AuthN() gin.HandlerFunc {
	return func(c *gin.Context) {
		door, err := c.Cookie(AuthNHTTPCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				c.Next()
				return
			} else {
				handler.HandleResponse(c, err, nil)
			}
			c.Abort()
			return
		}
		accessToken, err := m.ps.VerifyAccessToken(c.Request.Context(), door)
		switch {
		case err == nil && accessToken != nil:
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), ctxUserInfo, accessToken.GetUserInfo()))
		case errors.IsInternalServer(err):
			handler.HandleResponse(c, err, nil)
			c.Abort()
			return
		// user not login
		case err != nil || accessToken == nil:
			c.Next()
			return
		}
		c.Next()
	}
}

func (m *AuthNUserMiddleware) MustAuthN() gin.HandlerFunc {
	return func(c *gin.Context) {
		door, err := c.Cookie(AuthNHTTPCookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				handler.HandleResponse(c, errors.Unauthorized(reason.UnauthorizedError), nil)
			} else {
				handler.HandleResponse(c, err, nil)
			}
			c.Abort()
			return
		}
		accessToken, err := m.ps.VerifyAccessToken(c.Request.Context(), door)
		switch {
		case err == nil && accessToken != nil:
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), ctxUserInfo, accessToken.GetUserInfo()))
		case errors.IsInternalServer(err):
			handler.HandleResponse(c, err, nil)
			c.Abort()
			return
		case err != nil || accessToken == nil:
			handler.HandleResponse(c, errors.Unauthorized(reason.UnauthorizedError), nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func UserInfo(ctx context.Context) entity.UserInfo {
	userInfoValue := ctx.Value(ctxUserInfo)
	if userInfo, ok := (userInfoValue).(entity.UserInfo); ok {
		return userInfo
	}
	return entity.UserInfo{}
}
