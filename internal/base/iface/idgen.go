package iface

import "context"

type UniqueIDGenerator interface {
	NextIDInt64(ctx context.Context) (int64, error)
}
