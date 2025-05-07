package tracker

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	pbtr "github.com/Sp92535/proto/tracker/pb"
	"github.com/redis/go-redis/v9"
)

// ip limit to send
const LIMIT = 50

// implemented function to return resource holders of a particular file
func (tr *TrackerServer) GetResourceHolders(ctx context.Context, req *pbtr.GetResourceHoldersRequest) (*pbtr.GetResourceHoldersResponse, error) {

	var res pbtr.GetResourceHoldersResponse
	var currentSize int = 0

	res.Holder = make(map[uint64]*pbtr.HolderRow)

	// get ips for only the chunks that client doesn't have
	for _, chunkNo := range req.Need {
		if currentSize >= LIMIT {
			break
		}

		// construct key
		key := fmt.Sprintf("%s:%d", req.FileHash, chunkNo)

		// remove expired members before giving new ips
		RemoveExpired(key, ctx)

		// donot append if limit reached -> note: not exiting due to remaining initialization of ips
		sample := random(2)
		arr, err := Rdb.ZRandMember(ctx, key, sample).Result()
		if len(arr) > 0 {
			// append random number of ips from 1-5
			res.Holder[chunkNo] = &pbtr.HolderRow{}
			res.Holder[chunkNo].Ips = arr
			currentSize += len(res.Holder[chunkNo].Ips)
		}

		if err != nil {
			log.Printf("error fetching holders from redis: %v", err)
			return nil, err
		}
	}
	return &res, nil
}

// implemented function to register client's node server ip in chunk sets
func (tr *TrackerServer) RegisterResourceHolder(ctx context.Context, req *pbtr.RegisterResourceHolderRequest) (*pbtr.Empty, error) {

	msg := &SyncMessage{
		sender:    tr.address,
		node_ip:   req.Address,
		file_hash: req.FileHash,
	}

	for chunkNo, ok := range req.Status {
		if !ok {
			continue
		}
		msg.chunks = append(msg.chunks, uint64(chunkNo))
		// construct key
		key := fmt.Sprintf("%s:%d", req.FileHash, chunkNo)
		// construct score
		score := float64(time.Now().Add(TTL).Unix())

		// add to set along with score(expiration time)
		err := Rdb.ZAdd(ctx, key, redis.Z{Score: score, Member: req.Address}).Err()
		if err != nil {
			log.Printf("error adding holder to redis: %v", err)
			return nil, err
		}
	}
	if random(15) == random(15) {
		go Publish(msg)
	}

	return &pbtr.Empty{}, nil
}

// function to get random number
func random(x int) int {
	return rand.IntN(x) + 1
}
