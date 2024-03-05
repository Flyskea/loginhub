package server

import (
	"github.com/google/wire"

	"loginhub/internal/server/http"
)

var ProviderSetServer = wire.NewSet(
	http.ProviderSetHTTP,
)
