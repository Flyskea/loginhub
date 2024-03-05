package handler

import (
	"github.com/google/wire"

	"loginhub/internal/handler/http"
)

var ProviderSetHandler = wire.NewSet(
	http.ProviderSetHTTPHandler,
)
