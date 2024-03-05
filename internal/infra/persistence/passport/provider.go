package passport

import (
	"github.com/google/wire"

	"loginhub/internal/domain/passport/repository"
	"loginhub/internal/infra/persistence/passport/repo"
)

var ProviderSetPersistencePassport = wire.NewSet(
	repo.NewAccessTokenRepo,
	wire.Bind(new(repository.AccessTokenRepository), new(*repo.AccessTokenRepo)),
	repo.NewRefreshTokenRepo,
	wire.Bind(new(repository.RefreshTokenRepository), new(*repo.RefreshTokenRepo)),
	repo.NewLoginDeviceRepo,
	wire.Bind(new(repository.LoginDeviceRepository), new(*repo.LoginDeviceRepo)),
)
