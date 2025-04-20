package tracker

import (
	"context"
	"log"

	pbtr "github.com/Sp92535/proto/tracker/pb"
)

func (p *TrackerServer) GetResourceHolders(ctx context.Context, req *pbtr.GetResourceHoldersRequest) (*pbtr.GetResourceHoldersResponse, error) {

	ips, err := Rdb.SMembers(ctx, req.FileHash).Result()
	if err != nil {
		log.Printf("error fetching holders from redis: %v", err)
		return nil, err
	}

	return &pbtr.GetResourceHoldersResponse{
		Holders: ips,
	}, nil
}

func (p *TrackerServer) RegisterResourceHolder(ctx context.Context, req *pbtr.RegisterResourceHolderRequest) (*pbtr.Empty, error) {
	err := Rdb.SAdd(ctx, req.FileHash, req.Address).Err()
	if err != nil {
		log.Printf("error adding holder to redis: %v", err)
		return nil, err
	}
	return &pbtr.Empty{}, nil
}
