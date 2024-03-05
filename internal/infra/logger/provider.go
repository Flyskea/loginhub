package logger

import (
	"github.com/google/wire"

	"loginhub/pkg/logx"
)

var ProviderSetLogger = wire.NewSet(
	NewLogger,
	logx.NewKratosToSlog,
)
