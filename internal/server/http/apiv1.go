package http

import (
	"loginhub/internal/base/middleware/authn"
	"loginhub/internal/handler/http/passport"

	"github.com/gin-gonic/gin"
)

type APIV1Router struct {
	authnMiddleware *authn.AuthNUserMiddleware
	passportHandler *passport.PassportHandler
}

// @title        LoginHub API
// @version      v1
// @description  This is a auth server

// @BasePath  /api/v1
func NewAPIV1Router(
	authnMiddleware *authn.AuthNUserMiddleware,
	passportHandler *passport.PassportHandler,
) *APIV1Router {
	return &APIV1Router{
		authnMiddleware: authnMiddleware,
		passportHandler: passportHandler,
	}
}

func (r *APIV1Router) RegisterMustAuthHandler(g *gin.RouterGroup) {
	g.Use(r.authnMiddleware.MustAuthN())
	{
		passportAuthV1 := g.Group("/passport")
		{
			passportAuthV1.POST("/logout", r.passportHandler.Logout)
			passportAuthV1.GET("/device", r.passportHandler.LoginDevices)
			passportAuthV1.POST("/device/:id/kick", r.passportHandler.KickDevice)
		}
	}
}

func (r *APIV1Router) RegisterHandler(g *gin.RouterGroup) {
	{
		passportV1 := g.Group("/passport")
		{
			passportV1.POST("/session/refresh", r.passportHandler.RefreshCookie)
			passportV1.POST("/mail/send", r.passportHandler.EmailCaptchaSend)
			passportV1.POST("/register", r.passportHandler.Register)
			passportV1.POST("/login", r.passportHandler.Login)
		}
	}
}
