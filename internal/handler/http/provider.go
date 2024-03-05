package http

import (
	"github.com/google/wire"

	"loginhub/internal/handler/http/passport"
)

var ProviderSetHTTPHandler = wire.NewSet(
	passport.NewPassportHandler,
)
