package tracker

import (
	"context"

	pbtr "github.com/Sp92535/internal/tracker/pb"
)

func (p *TrackerServer) SendResourceRequest(ctx context.Context, req *pbtr.ResourceRequest) (*pbtr.Empty, error) {
	PublishMessage(req.Msg)
	return &pbtr.Empty{}, nil
}
