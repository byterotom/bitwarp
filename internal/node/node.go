package node

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Sp92535/pkg"
	"github.com/Sp92535/pkg/warpgen"
	pbtr "github.com/Sp92535/proto/tracker/pb"
)

type Node struct {
	address  string
	warp     *warpgen.Warp
	status   []bool
	holders  [][]string
	isSeeder bool
}

func NewNode(warp *warpgen.Warp, isSeeder bool) *Node {
	var n Node
	n = Node{
		address:  pkg.GetLocalIp(),
		warp:     warp,
		isSeeder: isSeeder,
		status:   make([]bool, warp.TotalChunks),
		holders:  make([][]string, warp.TotalChunks),
	}
	if isSeeder {
		for idx := range n.status {
			n.status[idx] = true
		}
	}
	return &n
}

func (n *Node) SendResourceRequest() {
	if n.isSeeder {
		return
	}
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
	for idx := range res.Holder {
		n.holders[idx] = res.Holder[idx].Ips
	}
	log.Println(n.holders)
}

func (n *Node) RegisterLoop() {
	ticker := time.NewTicker(30 * time.Second)

	n.Register()
	for {
		select {
		case <-ticker.C:
			n.Register()
		}
	}
}

func (n *Node) Register() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n.UpdateStatus()

	req := &pbtr.RegisterResourceHolderRequest{
		FileHash: n.warp.FileHash,
		Status:   n.status,
		Address:  n.address,
	}

	_, err := trackerClient.RegisterResourceHolder(ctx, req)
	if err != nil {
		log.Fatalf("could not invoke rpc: %v", err)
	}
}

func (n *Node) UpdateStatus() {
	if n.isSeeder {
		return
	}
	chunkDir := "storage/temp/" + n.warp.FileHash + "/"
	os.MkdirAll(chunkDir, os.ModePerm)

	// iterate over files
	files, err := os.ReadDir(chunkDir)
	if err != nil {
		log.Fatalf("error reading chunk directory %v", err)
	}

	for _, file := range files {
		filename := file.Name()
		if idx, err := strconv.Atoi(filename); err == nil {
			n.status[idx] = true
		}
	}
}
