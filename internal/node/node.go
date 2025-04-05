package node

import (
	"context"
	"log"
	"time"

	pbtr "github.com/Sp92535/internal/tracker/pb"
	"github.com/Sp92535/pkg"
	"github.com/Sp92535/pkg/warpgen"
)

type Node struct {
	ip     string
	warp   *warpgen.Warp
	status []bool
}

func NewNode(warp *warpgen.Warp) *Node {
	return &Node{
		ip:     pkg.GetLocalIp(),
		warp:   warp,
		status: make([]bool, warp.TotalChunks),
	}
}

func (n *Node) SendResourceRequest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := trackerClient.SendResourceRequest(ctx, &pbtr.ResourceRequest{FileHash: n.warp.FileHash})
	if err != nil {
		log.Fatalf("could not invoke rpc: %v", err)
	}
}
