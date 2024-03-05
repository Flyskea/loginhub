package http

import (
	"github.com/google/wire"

	"loginhub/internal/handler/http/oauth2"
	"loginhub/internal/handler/http/passport"
)

var ProviderSetHTTPHandler = wire.NewSet(
	oauth2.NewOAuth2ProviderHandler,
	passport.NewPassportHandler,
)
