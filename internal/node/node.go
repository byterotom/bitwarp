// node.go
package node

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Sp92535/pkg"
	"github.com/Sp92535/pkg/warpgen"
	pbno "github.com/Sp92535/proto/node/pb"
	pbtr "github.com/Sp92535/proto/tracker/pb"
)

type Node struct {
	address string
	warp    *warpgen.Warp

	mu        sync.Mutex
	available uint64
	status    []bool

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

func (n *Node) GetChunk(addr string, chunkNo int, once *sync.Once, wg *sync.WaitGroup) {
	defer wg.Done()

	nodeConn, nodeClient := NodeClient(addr)
	defer nodeConn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &pbno.GetResourceRequest{
		ChunkNo: uint64(chunkNo),
	}

	chunk, err := nodeClient.GetResource(ctx, req)
	if err != nil {
		log.Printf("could not invoke rpc: %v", err)
	}

	if warpgen.Verify(n.warp.Chunk[chunkNo], chunk.ChunkData) {
		once.Do(func() {
			n.mu.Lock()
			defer n.mu.Unlock()
			warpgen.CreateChunk(n.warp.FileHash, chunkNo, chunk.ChunkData)
			n.available++
			n.status[chunkNo] = true
			log.Printf("got chunk %d from %s.", chunkNo, addr)
		})
	}

}

func (n *Node) Download() {
	if n.isSeeder {
		return
	}

	for n.available != uint64(n.warp.TotalChunks) {

		var wg sync.WaitGroup

		n.SendResourceRequest()
		ips := n.holders
		for chunkNo := range ips {

			if n.status[chunkNo] {
				continue
			}

			var once sync.Once
			for _, addr := range ips[chunkNo] {
				wg.Add(1)
				go n.GetChunk(addr, chunkNo, &once, &wg)
			}
		}

		wg.Wait()
	}

	n.warp.MergeChunks()
}
