package repo

import (
	"github.com/google/wire"

	"loginhub/internal/base/iface"
	"loginhub/internal/infra/repo/captcha"
	"loginhub/internal/infra/repo/location"
	messagetemplate "loginhub/internal/infra/repo/message_template"
	"loginhub/internal/infra/repo/unique"
	servicecaptcha "loginhub/internal/service/captcha"
	servicemessagetemplate "loginhub/internal/service/message_template"
)

var ProviderSetInfraRepo = wire.NewSet(
	messagetemplate.NewMessageTemplateRepo,
	wire.Bind(new(servicemessagetemplate.MessageTemplateRepository), new(*messagetemplate.MessageTemplateRepo)),
	unique.NewSnowflake,
	unique.NewUniqueIDRepo,
	wire.Bind(new(iface.UniqueIDGenerator), new(*unique.UniqueIDRepo)),
	location.NewIPRegionSearcher,
	wire.Bind(new(iface.IP2RegionSearcher), new(*location.XDBIPRegionSearcher)),
	captcha.NewCaptchaRepo,
	wire.Bind(new(servicecaptcha.CaptchaRepository), new(*captcha.CaptchaRepo)),
)
