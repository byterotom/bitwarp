package node

import (
	"context"
	"time"

	pbno "github.com/Sp92535/proto/node/pb"
)

func (nodeServer *NodeServer) Ping(ctx context.Context, req *pbno.Empty) (*pbno.Pong, error) {
	return &pbno.Pong{
		Time: time.Now().String(),
	}, nil
}

func (nodeServer *NodeServer) GetResource(ctx context.Context, req *pbno.GetResourceRequest) (*pbno.GetResourceResponse, error) {

	if !nodeServer.isSeeder {
		return &pbno.GetResourceResponse{
			ChunkNo:   req.ChunkNo,
			ChunkData: []byte{},
		}, nil
	}

	return &pbno.GetResourceResponse{
		ChunkNo:   req.ChunkNo,
		ChunkData: nodeServer.warp.ReadChunk(int(req.ChunkNo)),
	}, nil
}
