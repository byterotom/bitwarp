package node

import (
	"context"
	"log"
	"os"

	pbtr "github.com/Sp92535/internal/tracker/pb"
	"github.com/Sp92535/pkg"
	"github.com/Sp92535/pkg/warpgen"
)

type Node struct {
	address string
	warp    *warpgen.Warp
	status  []bool
	holders []string
}

func NewNode(warp *warpgen.Warp) *Node {
	return &Node{
		address: pkg.GetLocalIp(),
		warp:    warp,
		status:  make([]bool, warp.TotalChunks),
		holders: make([]string, 20),
	}
}

func (n *Node) SendResourceRequest() {

	chunkDir := "storage/temp/" + n.warp.FileHash + "/"
	os.MkdirAll(chunkDir, os.ModePerm)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &pbtr.GetResourceHoldersRequest{
		FileHash: n.warp.FileHash,
	}
	res, err := trackerClient.GetResourceHolders(ctx, req)
	if err != nil {
		log.Fatalf("could not invoke rpc: %v", err)
	}
	n.holders = res.Holders
	log.Println(n.holders)
}

func (n *Node) SendResourceResponse(fileHash string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if n.warp.FileHash != fileHash {
		return
	}

	req := &pbtr.RegisterResourceHolderRequest{
		FileHash: n.warp.FileHash,
		Address:  n.address,
	}

	_, err := trackerClient.RegisterResourceHolder(ctx, req)
	if err != nil {
		log.Fatalf("could not invoke rpc: %v", err)
	}
}
