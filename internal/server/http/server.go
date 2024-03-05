package http

import (
	transportgin "github.com/tx7do/kratos-transport/transport/gin"

	"loginhub/internal/base/middleware/lang"
	"loginhub/internal/conf"
	"loginhub/pkg/logx"
)

func NewHTTPServer(
	conf *conf.HTTP,
	apiv1Router *APIV1Router,
	logger *logx.KratosToSlog,
) *transportgin.Server {
	server := transportgin.NewServer(
		transportgin.WithAddress(conf.Addr),
		transportgin.WithTimeout(conf.Timeout.AsDuration()),
	)

	server.Use(lang.ExtractAndSetAcceptLanguage, transportgin.GinRecovery(logger, true), transportgin.GinLogger(logger))
	authv1 := server.Group("/api/v1")
	apiv1Router.RegisterMustAuthV1APIRouter(authv1)
	unauthv1 := server.Group("/api/v1")
	apiv1Router.RegisterNoAuthV1APIRouter(unauthv1)

	return server
}
