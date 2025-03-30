package node

import (
	"context"
	"time"

	pbno "github.com/Sp92535/internal/node/pb"
)

func (p *NodeServer) Ping(ctx context.Context, req *pbno.Empty) (*pbno.Pong, error) {
	return &pbno.Pong{
		Time: time.Now().String(),
	}, nil
}
