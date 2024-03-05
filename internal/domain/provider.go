package domain

import (
	"github.com/google/wire"

	"loginhub/internal/domain/passport/service"
)

var ProviderSetDomain = wire.NewSet(
	service.NewPassportService,
)
