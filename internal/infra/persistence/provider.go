package persistence

import (
	"github.com/google/wire"

	"loginhub/internal/domain/user/repository"
	"loginhub/internal/infra/persistence/passport"
	"loginhub/internal/infra/persistence/user"
)

var ProviderSetPersistence = wire.NewSet(
	passport.ProviderSetPersistencePassport,
	user.New,
	wire.Bind(new(repository.UserRepository), new(*user.UserRepo)),
)
