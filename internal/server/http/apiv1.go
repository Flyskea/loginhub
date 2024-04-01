package http

import (
	"github.com/gin-gonic/gin"

	"loginhub/internal/base/middleware/authn"
	"loginhub/internal/handler/http/oauth2"
	"loginhub/internal/handler/http/passport"
)

type APIV1Router struct {
	authnMiddleware *authn.AuthNUserMiddleware
	oauth2Handler   *oauth2.OAuth2Handler
	passportHandler *passport.PassportHandler
}

// @title        LoginHub API
// @version      v1
// @description  This is a auth server

// @BasePath  /api/v1
func NewAPIV1Router(
	authnMiddleware *authn.AuthNUserMiddleware,
	oauth2Handler *oauth2.OAuth2Handler,
	passportHandler *passport.PassportHandler,
) *APIV1Router {
	return &APIV1Router{
		authnMiddleware: authnMiddleware,
		oauth2Handler:   oauth2Handler,
		passportHandler: passportHandler,
	}
}

func (r *APIV1Router) RegisterMustAuthV1APIRouter(g *gin.RouterGroup) {
	g.Use(r.authnMiddleware.MustAuthN())
	{
		passportAuthV1 := g.Group("/passport")
		{
			passportAuthV1.POST("/logout", r.passportHandler.Logout)
			passportAuthV1.GET("/device", r.passportHandler.LoginDevices)
			passportAuthV1.POST("/device/:id/kick", r.passportHandler.KickDevice)
		}
	}
	{
		oauth2AuthV1 := g.Group("/oauth2")
		{
			oauth2AuthV1.POST("/provider", r.oauth2Handler.CreateProvider)
			oauth2AuthV1.GET("/providers", r.oauth2Handler.ListProviders)
			oauth2AuthV1.POST("/provider/:id/update", r.oauth2Handler.UpdateProvider)
			oauth2AuthV1.POST("/provider/:id/delete", r.oauth2Handler.DeleteProvider)
		}
	}
}

func (r *APIV1Router) RegisterNoAuthV1APIRouter(g *gin.RouterGroup) {
	{
		passportV1 := g.Group("/passport")
		{
			passportV1.POST("/session/refresh", r.passportHandler.RefreshCookie)
			passportV1.POST("/mail/send", r.passportHandler.EmailCaptchaSend)
			passportV1.POST("/register", r.passportHandler.Register)
			passportV1.POST("/login", r.passportHandler.Login)
		}
	}
	{
		oauth2V1 := g.Group("/oauth2")
		oauth2V1.GET("/provider/supported", r.oauth2Handler.OAuth2SupportedProviders)
		oauth2V1.GET("/redirect/:provider", r.oauth2Handler.OAuth2RequestURL)
	}
}
