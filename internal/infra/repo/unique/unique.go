package unique

import (
	"context"

	"github.com/bwmarrin/snowflake"

	"loginhub/internal/base/iface"
)

var _ iface.UniqueIDGenerator = (*UniqueIDRepo)(nil)

type UniqueIDRepo struct {
	node *snowflake.Node
}

func NewSnowflake() *snowflake.Node {
	node, err :=  snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return node
}

func NewUniqueIDRepo(node *snowflake.Node) *UniqueIDRepo {
	return &UniqueIDRepo{node: node}
}

func (r *UniqueIDRepo) NextIDInt64(ctx context.Context) (int64, error) {
	return r.node.Generate().Int64(), nil
}
