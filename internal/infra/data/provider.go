package data

import (
	"github.com/google/wire"

	"loginhub/internal/base/iface"
)

var ProviderSetData = wire.NewSet(
	NewDB,
	NewTXManager,
	wire.Bind(new(iface.Transaction), new(*TXManager)),
	NewRDB,
)
