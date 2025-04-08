package tracker

import (
	"context"
	"sync"

	pbtr "github.com/Sp92535/internal/tracker/pb"
)

var mu sync.Mutex
var db map[string]map[string]struct{} = make(map[string]map[string]struct{})

func (p *TrackerServer) GetResourceHolders(ctx context.Context, req *pbtr.GetResourceHoldersRequest) (*pbtr.GetResourceHoldersResponse, error) {

	PublishRequest(req.FileHash)
	var holders []string
	for ip := range db[req.FileHash] {
		holders = append(holders, ip)
	}

	return &pbtr.GetResourceHoldersResponse{
		Holders: holders,
	}, nil
}

func (p *TrackerServer) RegisterResourceHolder(ctx context.Context, res *pbtr.RegisterResourceHolderRequest) (*pbtr.Empty, error) {
	mu.Lock()
	defer mu.Unlock()
	if db[res.FileHash] == nil {
		db[res.FileHash] = make(map[string]struct{})
	}
	db[res.FileHash][res.Address] = struct{}{}

	return &pbtr.Empty{}, nil
}
