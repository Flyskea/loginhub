package oauth2

import (
	"github.com/google/wire"

	"loginhub/internal/domain/oauth2/repository"
	"loginhub/internal/infra/persistence/oauth2/repo"
)

var ProviderSetPersistenceOAuth2 = wire.NewSet(
	repo.NewOAuth2ProviderRepo,
	wire.Bind(new(repository.OAuth2ProviderInfoRepository), new(*repo.OAuth2ProviderRepo)),
	repo.NewUserThirdAuthRepo,
	wire.Bind(new(repository.OAuth2UserRepository), new(*repo.UserThirdAuthRepo)),
)
