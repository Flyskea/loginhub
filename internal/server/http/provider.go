package http

import (
	"loginhub/internal/base/middleware/authn"

	"github.com/google/wire"
)

var ProviderSetHTTP = wire.NewSet(
	authn.NewAuthNUserMiddleware,
	NewAPIV1Router,
	NewHTTPServer,
)
