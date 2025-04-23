package node

import (
	"context"
	"time"

	pbno "github.com/Sp92535/proto/node/pb"
)

// implemented ping function to test connectivity
func (nodeServer *NodeServer) Ping(ctx context.Context, req *pbno.Empty) (*pbno.Pong, error) {
	return &pbno.Pong{
		Time: time.Now().String(),
	}, nil
}

// implemented function to get resource
func (nodeServer *NodeServer) GetResource(ctx context.Context, req *pbno.GetResourceRequest) (*pbno.GetResourceResponse, error) {
	
	data, err := nodeServer.warp.ReadChunk(int(req.ChunkNo), nodeServer.isSeeder)
	return &pbno.GetResourceResponse{
		ChunkNo:   req.ChunkNo,
		ChunkData: data,
	}, err
}
