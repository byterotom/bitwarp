// node.go
package node

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Sp92535/pkg"
	"github.com/Sp92535/pkg/warpgen"
	pbno "github.com/Sp92535/proto/node/pb"
	pbtr "github.com/Sp92535/proto/tracker/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Node struct {
	address string
	warp    *warpgen.Warp

	mu     sync.Mutex
	status []bool
	need   []uint64

	holders     map[uint64][]string
	holdersConn map[string](*grpc.ClientConn)

	isSeeder bool
}

// constructor function to create new node
func NewNode(warp *warpgen.Warp, isSeeder bool) *Node {
	var n Node
	n = Node{
		address:     pkg.GetLocalIp(),
		warp:        warp,
		isSeeder:    isSeeder,
		status:      make([]bool, warp.TotalChunks),
		need:        make([]uint64, warp.TotalChunks),
		holders:     make(map[uint64][]string),
		holdersConn: make(map[string]*grpc.ClientConn),
	}
	// if is seeder set all status to true
	if isSeeder {
		for idx := range n.status {
			n.status[idx] = true
		}
	} else {
		for idx := range n.need {
			n.need[idx] = uint64(idx)
		}
	}
	return &n
}

const LOOP_INTERVAL = 50 * time.Second

// function to loop register having resources to client and renew
func (n *Node) RegisterLoop() {
	ticker := time.NewTicker(LOOP_INTERVAL)

	// initial register
	n.register()
	for {
		select {
		case <-ticker.C:
			n.register()
		}
	}
}

// function to register resource to tracker
func (n *Node) register() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &pbtr.RegisterResourceHolderRequest{
		FileHash: n.warp.FileHash,
		Status:   n.status,
		Address:  n.address,
	}

	_, err := trackerClient.RegisterResourceHolder(ctx, req)
	if err != nil {
		log.Printf("could not invoke rpc: %v", err)
	}
}

// function to send holder request to tracker
func (n *Node) sendHolderRequest() {
	if n.isSeeder {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &pbtr.GetResourceHoldersRequest{
		FileHash: n.warp.FileHash,
		Need:     pkg.GetRandom(n.need, 50),
	}

	res, err := trackerClient.GetResourceHolders(ctx, req)
	if err != nil {
		log.Printf("could not invoke rpc: %v", err)
		return
	}
	for key := range res.Holder {
		n.holders[key] = res.Holder[key].Ips
		for _, addr := range n.holders[key] {
			if _, ok := n.holdersConn[addr]; !ok {
				n.holdersConn[addr] = NodeConn(addr)
			}
		}
	}

	// log.Println(n.holders)
}

// function to close all connections
func (n *Node) closeAllConn() {
	for addr, conn := range n.holdersConn {
		conn.Close()
		delete(n.holdersConn, addr)
	}
}

// function to fetch chunk from a node
func (n *Node) getChunk(ctx context.Context, cancel context.CancelFunc, addr string, chunkNo uint64, once *sync.Once, wg *sync.WaitGroup) {
	defer wg.Done()

	// initialize node client
	nodeClient := pbno.NewNodeServiceClient(n.holdersConn[addr])

	req := &pbno.GetResourceRequest{
		ChunkNo: uint64(chunkNo),
	}

	chunk, err := nodeClient.GetResource(ctx, req)
	if err != nil {
		if status.Code(err) != codes.Canceled {
			log.Printf("could not invoke rpc: %v", err)
		}
		return
	}

	// verify the chunk
	if warpgen.Verify(n.warp.Chunk[chunkNo], chunk.ChunkData) {
		// save once
		once.Do(func() {
			n.mu.Lock()
			defer n.mu.Unlock()
			warpgen.CreateChunk(n.warp.FileHash, chunkNo, chunk.ChunkData)
			n.status[chunkNo] = true
			n.need = pkg.Remove(n.need, uint64(chunkNo))
			log.Printf("got chunk %d from %s", chunkNo, addr)
			cancel()
		})
	}
}

const TIMEOUT = 10 * time.Second
const WORKERS = 50

// function to download the file
func (n *Node) Download() {
	// dont download if seeder
	if n.isSeeder {
		return
	}

	// loop until all chunks are found
	for len(n.need) > 0 {
		// waitgroup to wait until response
		var wg sync.WaitGroup
		worker := make(chan struct{}, WORKERS)

		n.sendHolderRequest()
		ips := n.holders
		for chunkNo := range ips {

			if n.status[chunkNo] {
				continue
			}

			// sync once to get the fastest
			ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)

			var once sync.Once
			var innerWg sync.WaitGroup

			for _, addr := range ips[chunkNo] {
				innerWg.Add(1)
				go n.getChunk(ctx, cancel, addr, chunkNo, &once, &innerWg)
			}

			worker <- struct{}{}
			wg.Add(1)
			go func() {
				defer func() {
					<-worker
					wg.Done()
					cancel()
				}()
				innerWg.Wait()
			}()

		}
		wg.Wait()
		close(worker)
		n.closeAllConn()
		n.register()
	}
}

// function to update chunks status in case of failure
func (n *Node) UpdateStatus() bool {
	if n.isSeeder {
		return true
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	chunkDir := "storage/temp/" + n.warp.FileHash + "/"
	os.MkdirAll(chunkDir, os.ModePerm)

	// iterate over chunk files
	files, err := os.ReadDir(chunkDir)
	if err != nil {
		log.Printf("error reading chunk directory %v", err)
	}

	var flag bool = true
	for _, file := range files {
		filename := file.Name()
		if idx, err := strconv.Atoi(filename); err == nil && !n.status[idx] {
			n.status[idx] = true
			n.need = pkg.Remove(n.need, uint64(idx))
			flag = false
		}
	}

	return flag
}
