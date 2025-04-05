package tracker

import (
	"context"

	pbtr "github.com/Sp92535/internal/tracker/pb"
)

// function to send resource request to all peers via tracker
func (p *TrackerServer) SendResourceRequest(ctx context.Context, req *pbtr.ResourceRequest) (*pbtr.Empty, error) {
	PublishRequest(req.FileHash)
	return &pbtr.Empty{}, nil
}
