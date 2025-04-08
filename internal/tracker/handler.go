package tracker

import (
	"context"

	pbtr "github.com/Sp92535/internal/tracker/pb"
)

var db map[string][][]string

func (p *TrackerServer) GetResourceHolders(ctx context.Context, req *pbtr.GetResourceHoldersRequest) (*pbtr.GetResourceHoldersResponse, error) {
	if db == nil {
		db = make(map[string][][]string)
	}

	db[req.FileHash] = make([][]string, len(req.Status))
	PublishRequest(req.FileHash)
	return &pbtr.GetResourceHoldersResponse{
		FileHash: req.FileHash,
		Holders:  []string{"a", "b", "c"},
	}, nil
}

func (p *TrackerServer) RegisterResourceHolder(ctx context.Context, res *pbtr.RegisterResourceHolderRequest) (*pbtr.Empty, error) {
	if db == nil {
		db = make(map[string][][]string)
	}

	adr, ok := db[res.FileHash]
	if ok {
		for idx, val := range res.Status {
			if val {
				adr[idx] = append(adr[idx], res.Address)
			}
		}
	}

	return &pbtr.Empty{}, nil
}
