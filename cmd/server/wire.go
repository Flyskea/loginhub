//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"

	"loginhub/internal/conf"
	"loginhub/internal/domain"
	"loginhub/internal/handler"
	"loginhub/internal/infra/data"
	"loginhub/internal/infra/email"
	"loginhub/internal/infra/logger"
	"loginhub/internal/infra/persistence"
	"loginhub/internal/infra/repo"
	"loginhub/internal/server"
	"loginhub/internal/service"
)

func initApp(
	logConf *conf.Log,
	dbConf *conf.Database,
	rdbConf *conf.Redis,
	passportConf *conf.Passport,
	captchaConf *conf.Captcha,
	smtpConf *conf.SMTP,
	httpConf *conf.HTTP,
	ip2RegionConf *conf.IP2Region,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		logger.ProviderSetLogger,
		data.ProviderSetData,
		email.NewEmailSender,
		persistence.ProviderSetPersistence,
		repo.ProviderSetInfraRepo,
		domain.ProviderSetDomain,
		handler.ProviderSetHandler,
		service.ProviderSetService,
		server.ProviderSetServer,
		newApp,
	))
}
