package tracker

import (
	"context"
	"fmt"
	"log"

	pbtr "github.com/Sp92535/proto/tracker/pb"
)

func (p *TrackerServer) GetResourceHolders(ctx context.Context, req *pbtr.GetResourceHoldersRequest) (*pbtr.GetResourceHoldersResponse, error) {

	var res pbtr.GetResourceHoldersResponse
	var err error

	res.Holder = make([]*pbtr.HolderRow, len(req.Status))
	for chunkNo, ok := range req.Status {
		if ok {
			continue
		}
		key := fmt.Sprintf("%s:%d", req.FileHash, chunkNo)

		res.Holder[chunkNo] = &pbtr.HolderRow{}
		res.Holder[chunkNo].Ips, err = Rdb.SMembers(ctx, key).Result()
		if err != nil {
			log.Printf("error fetching holders from redis: %v", err)
			return nil, err
		}
	}
	return &res, nil
}

func (p *TrackerServer) RegisterResourceHolder(ctx context.Context, req *pbtr.RegisterResourceHolderRequest) (*pbtr.Empty, error) {
	for chunkNo, ok := range req.Status {
		if !ok {
			continue
		}

		key := fmt.Sprintf("%s:%d", req.FileHash, chunkNo)

		err := Rdb.SAdd(ctx, key, req.Address).Err()
		if err != nil {
			log.Printf("error adding holder to redis: %v", err)
			return nil, err
		}
	}
	return &pbtr.Empty{}, nil
}
