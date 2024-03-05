package service

import (
	"github.com/google/wire"

	"loginhub/internal/service/captcha"
	messagetemplate "loginhub/internal/service/message_template"
	"loginhub/internal/service/passport"
)

var ProviderSetService = wire.NewSet(
	captcha.New,
	messagetemplate.NewMessageTemplateService,
	passport.NewPassportService,
)
