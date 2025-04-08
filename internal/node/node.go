package node

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"os"
	"strconv"

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
		Status:   n.status,
	}
	res, err := trackerClient.GetResourceHolders(ctx, req)
	if err != nil {
		log.Fatalf("could not invoke rpc: %v", err)
	}
	n.holders = res.Holders
}

func (n *Node) SendResourceResponse(fileHash string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chunkDir := "storage/temp/" + fileHash + "/"

	chunks, err := os.ReadDir(chunkDir)
	if errors.Is(err, fs.ErrNotExist) {
		return
	}
	if err != nil {
		log.Fatalf("error reading chunks dir: %v", err)
	}

	for _, chunkNo := range chunks {
		idx, _ := strconv.Atoi(chunkNo.Name())
		n.status[idx] = true
	}

	req := &pbtr.RegisterResourceHolderRequest{
		FileHash: n.warp.FileHash,
		Status:   n.status,
		Address:  n.address,
	}

	_, err = trackerClient.RegisterResourceHolder(ctx, req)
	if err != nil {
		log.Fatalf("could not invoke rpc: %v", err)
	}
}
