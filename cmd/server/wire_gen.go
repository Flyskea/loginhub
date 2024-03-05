// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"loginhub/internal/base/middleware/authn"
	"loginhub/internal/conf"
	"loginhub/internal/domain/passport/service"
	oauth2_2 "loginhub/internal/handler/http/oauth2"
	passport2 "loginhub/internal/handler/http/passport"
	"loginhub/internal/infra/data"
	"loginhub/internal/infra/email"
	"loginhub/internal/infra/logger"
	repo2 "loginhub/internal/infra/persistence/oauth2/repo"
	"loginhub/internal/infra/persistence/passport/repo"
	"loginhub/internal/infra/persistence/user"
	"loginhub/internal/infra/repo/captcha"
	"loginhub/internal/infra/repo/location"
	"loginhub/internal/infra/repo/message_template"
	"loginhub/internal/infra/repo/unique"
	"loginhub/internal/server/http"
	captcha2 "loginhub/internal/service/captcha"
	messagetemplate2 "loginhub/internal/service/message_template"
	"loginhub/internal/service/oauth2"
	"loginhub/internal/service/passport"
	"loginhub/pkg/logx"
)

// Injectors from wire.go:

func initApp(logConf *conf.Log, dbConf *conf.Database, rdbConf *conf.Redis, passportConf *conf.Passport, captchaConf *conf.Captcha, smtpConf *conf.SMTP, httpConf *conf.HTTP, ip2RegionConf *conf.IP2Region) (*kratos.App, func(), error) {
	slogLogger := logger.NewLogger(logConf)
	kratosToSlog := logx.NewKratosToSlog(slogLogger)
	db, err := data.NewDB(dbConf)
	if err != nil {
		return nil, nil, err
	}
	txManager := data.NewTXManager(db)
	userRepo := user.New(txManager, slogLogger)
	passportService := service.NewPassportService(userRepo)
	universalClient, cleanup, err := data.NewRDB(rdbConf)
	if err != nil {
		return nil, nil, err
	}
	xdbipRegionSearcher, err := location.NewIPRegionSearcher(ip2RegionConf)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	loginDeviceRepo := repo.NewLoginDeviceRepo(txManager, universalClient, xdbipRegionSearcher)
	accessTokenRepo := repo.NewAccessTokenRepo(db, universalClient)
	refreshTokenRepo := repo.NewRefreshTokenRepo(universalClient, slogLogger)
	userThirdAuthRepo := repo2.NewUserThirdAuthRepo(txManager)
	node := unique.NewSnowflake()
	uniqueIDRepo := unique.NewUniqueIDRepo(node)
	oAuth2ProviderRepo := repo2.NewOAuth2ProviderRepo(txManager, universalClient)
	oAuth2Service := oauth2.NewOauth2Service(oAuth2ProviderRepo)
	captchaRepo := captcha.NewCaptchaRepo(universalClient)
	emailSender, cleanup2, err := email.NewEmailSender(smtpConf, slogLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	messageTemplateRepo := messagetemplate.NewMessageTemplateRepo(txManager)
	messageTemplateService := messagetemplate2.NewMessageTemplateService(messageTemplateRepo)
	captchaService := captcha2.New(captchaRepo, emailSender, messageTemplateService)
	passportPassportService, err := passport.NewPassportService(txManager, passportService, loginDeviceRepo, accessTokenRepo, refreshTokenRepo, userThirdAuthRepo, userRepo, uniqueIDRepo, oAuth2Service, captchaService)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	authNUserMiddleware := authn.NewAuthNUserMiddleware(passportPassportService)
	oAuth2Handler := oauth2_2.NewOAuth2ProviderHandler(oAuth2Service)
	passportHandler, err := passport2.NewPassportHandler(passportPassportService, xdbipRegionSearcher)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	apiv1Router := http.NewAPIV1Router(authNUserMiddleware, oAuth2Handler, passportHandler)
	server := http.NewHTTPServer(httpConf, apiv1Router, kratosToSlog)
	app := newApp(kratosToSlog, server)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
