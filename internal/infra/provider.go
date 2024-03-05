package infra

import (
	"github.com/google/wire"

	"loginhub/internal/infra/data"
	"loginhub/internal/infra/email"
	"loginhub/internal/infra/repo"
)

var ProviderSetInfra = wire.NewSet(
	data.ProviderSetData,
	email.NewEmailSender,
	repo.ProviderSetInfraRepo,
)
